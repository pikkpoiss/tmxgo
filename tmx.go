// Copyright 2014 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tmxgo

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
)

// The tilewidth and tileheight properties determine the general grid
// size of the map. The individual tiles may have different sizes.
// Larger tiles will extend at the top and right (anchored to the bottom left).
type Map struct {
	// The TMX format version, generally 1.0.
	Version string `xml:"version,attr"`

	// Map orientation. Tiled supports "orthogonal", "isometric"
	// and "staggered" (since 0.9.0) at the moment.
	Orientation string `xml:"orientation,attr"`

	// The map width in tiles.
	Width int32 `xml:"width,attr"`

	// The map height in tiles.
	Height int32 `xml:"height,attr"`

	// The width of a tile.
	TileWidth int32 `xml:"tilewidth,attr"`

	// The height of a tile.
	TileHeight int32 `xml:"tileheight,attr"`

	// The background color of the map. (since 0.9.0).
	BackgroundColor string `xml:"backgroundcolor,attr"`

	// Can contain properties.
	Properties []Property `xml:"properties>property"`

	// Can contain tileset.
	Tilesets []Tileset `xml:"tileset"`

	// Can contain layer.
	Layers []Layer `xml:"layer"`

	// Can contain objectgroup.
	ObjectGroups []ObjectGroup `xml:"objectgroup"`

	// Can contain imagelayer.
	ImageLayers []ImageLayer `xml:"imagelayer"`
}

func (m *Map) TilesFromLayerName(name string) (t []*Tile, err error) {
	for i := 0; i < len(m.Layers); i++ {
		if m.Layers[i].Name == name {
			return m.tilesFromLayer(&m.Layers[i])
		}
	}
	err = fmt.Errorf("No layer with name %v", name)
	return
}

func (m *Map) TilesFromLayerIndex(index int32) (t []*Tile, err error) {
	if index < 0 || index > int32(len(m.Layers)) {
		err = fmt.Errorf("Index %v out of bounds", index)
		return
	}
	return m.tilesFromLayer(&m.Layers[index])
}

func (m *Map) tilesFromLayer(layer *Layer) (t []*Tile, err error) {
	var datatiles []DataTile
	if datatiles, err = layer.Data.Tiles(); err != nil {
		return
	}
	sort.Sort(byFirstGid(m.Tilesets)) // Should be sorted but just in case.
	t = make([]*Tile, len(datatiles))
	for i := 0; i < len(datatiles); i++ {
		var (
			mapheight  = layer.Height * m.TileHeight
			tilebounds = Bounds{
				Y: float32(mapheight - ((int32(i) / layer.Width) * m.TileHeight)),
				X: float32((int32(i) % layer.Width) * m.TileWidth),
				W: float32(m.TileWidth),
				H: float32(m.TileHeight),
			}
			gid = datatiles[i].Gid
		)

		if gid == 0 {
			t[i] = nil
		} else if t[i], err = newTile(gid, m.Tilesets, tilebounds); err != nil {
			return
		}
	}
	return
}

type Bounds struct {
	X, Y, W, H float32
}

type Tile struct {
	Index         uint32
	Tileset       *Tileset
	FlipVert      bool
	FlipHorz      bool
	FlipDiag      bool
	TileBounds    Bounds
	TextureBounds Bounds
}

func (t *Tile) Triangles() [30]float32 {
	return [30]float32{
		t.TileBounds.X, t.TileBounds.Y, 0.0,
		t.TextureBounds.X, t.TextureBounds.Y,

		t.TileBounds.X + t.TileBounds.H, t.TileBounds.Y + t.TileBounds.H, 0.0,
		t.TextureBounds.X + t.TextureBounds.H, t.TextureBounds.Y + t.TextureBounds.H,

		t.TileBounds.X, t.TileBounds.Y + t.TileBounds.H, 0.0,
		t.TextureBounds.X, t.TextureBounds.Y + t.TextureBounds.H,

		t.TileBounds.X, t.TileBounds.Y, 0.0,
		t.TextureBounds.X, t.TextureBounds.Y,

		t.TileBounds.X + t.TileBounds.H, t.TileBounds.Y, 0.0,
		t.TextureBounds.X + t.TextureBounds.H, t.TextureBounds.Y,

		t.TileBounds.X + t.TileBounds.H, t.TileBounds.Y + t.TileBounds.H, 0.0,
		t.TextureBounds.X + t.TextureBounds.H, t.TextureBounds.Y + t.TextureBounds.H,
	}
}

const (
	FLIPPED_H_FLAG uint32 = 0x80000000
	FLIPPED_V_FLAG uint32 = 0x40000000
	FLIPPED_D_FLAG uint32 = 0x20000000
	CLEAR_FLIP     uint32 = (FLIPPED_H_FLAG | FLIPPED_V_FLAG | FLIPPED_D_FLAG)
)

func parseGid(gid uint32) (id uint32, fliph, flipv, flipd bool) {
	fliph = (gid & FLIPPED_H_FLAG) > 0
	flipv = (gid & FLIPPED_V_FLAG) > 0
	flipd = (gid & FLIPPED_D_FLAG) > 0
	id = gid & ^CLEAR_FLIP
	return
}

// The tilesets argument must first be sorted by firstgid.
func newTile(gid uint32, tilesets []Tileset, tilebounds Bounds) (t *Tile, err error) {
	var (
		tileset *Tileset
		count   = len(tilesets)
		fliph   bool
		flipv   bool
		flipd   bool
		index   uint32
	)
	if count == 0 {
		err = fmt.Errorf("No tilesets")
		return
	}
	gid, fliph, flipv, flipd = parseGid(gid)
	for i := 1; i < count; i++ {
		if gid < tilesets[i].FirstGid {
			tileset = &tilesets[i-1]
			break
		}
	}
	if tileset == nil {
		tileset = &tilesets[count-1]
	}
	index = gid - tileset.FirstGid
	t = &Tile{
		Index:         index,
		Tileset:       tileset,
		FlipVert:      flipv,
		FlipHorz:      fliph,
		FlipDiag:      flipd,
		TileBounds:    tilebounds,
		TextureBounds: tileset.TextureBounds(index),
	}
	return
}

// Sorts Tilesets by FirstGid property.
type byFirstGid []Tileset

func (b byFirstGid) Len() int           { return len(b) }
func (b byFirstGid) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byFirstGid) Less(i, j int) bool { return b[i].FirstGid < b[j].FirstGid }

type Tileset struct {
	// The first global tile ID of this tileset.
	// (this global ID maps to the first tile in this tileset).
	FirstGid uint32 `xml:"firstgid,attr"`

	// If this tileset is stored in an external TSX (Tile Set XML) file,
	// this attribute refers to that file. That TSX file has the
	// same structure as the attribute as described here.
	// (There is the firstgid attribute missing and this source
	// attribute is also not there. These two attributes are kept
	// in the TMX map, since they are map specific.)
	Source string `xml:"source,attr"`

	// The name of this tileset.
	Name string `xml:"name,attr"`

	// The (maximum) width of the tiles in this tileset.
	TileWidth int32 `xml:"tilewidth,attr"`

	// The (maximum) height of the tiles in this tileset.
	TileHeight int32 `xml:"tileheight,attr"`

	// The spacing in pixels between the tiles in this tileset.
	// (applies to the tileset image).
	Spacing int32 `xml:"spacing,attr"`

	// The margin around the tiles in this tileset.
	// (applies to the tileset image).
	Margin int32 `xml:"margin,attr"`

	// Can contain tileoffset (since 0.8.0).
	TileOffset *TileOffset `xml:"tileoffset"`

	// Can contain properties (since 0.8.0).
	Properties []Property `xml:"properties>property"`

	// Can contain image.
	Image *Image `xml:"image"`

	// Can contain terraintypes (since 0.9.0).
	TerrainTypes []Terrain `xml:"terraintypes>terrain"`

	// Can contain tile.
	TilesetTile []TilesetTile `xml:"tile"`
}

func (t *Tileset) TextureBounds(index uint32) Bounds {
	if t.Image == nil {
		return Bounds{0, 0, 0, 0}
	}
	var (
		tileswide = t.Image.Width / t.TileWidth
		tileshigh = t.Image.Height / t.TileHeight
	)
	return Bounds{
		Y: float32((t.Image.Height-int32(index)/tileswide)*t.TileHeight) / float32(t.Image.Height),
		X: float32((int32(index)%tileshigh)*t.TileWidth) / float32(t.Image.Width),
		W: float32(t.TileWidth) / float32(t.Image.Width),
		H: -float32(t.TileHeight) / float32(t.Image.Height),
	}
}

// This element is used to specify an offset in pixels,
// to be applied when drawing a tile from the related tileset.
// When not present, no offset is applied.
type TileOffset struct {
	// Horizontal offset in pixels.
	X int32 `xml:"x,attr"`

	// Vertical offset in pixels (positive is down).
	Y int32 `xml:"y,attr"`
}

// As of the current version of Tiled Qt, each tileset hass a single image
// associated with it, which is cut into smaller tiles based on the
// attributes defined on the tileset element. Later versions may
// add support for adding multiple images to a single tileset,
// as is possible in Tiled Java.
type Image struct {
	// Used for embedded images, in combination with a data child element.
	// (since 0.9.0)
	Format string `xml:"format,attr"`

	// Used by some versions of Tiled Java.
	// Deprecated and unsupported by Tiled Qt.
	Id int32 `xml:"id,attr"`

	// The reference to the tileset image file.
	// (Tiled supports most common image formats).
	Source string `xml:"source,attr"`

	// Defines a specific color that is treated as transparent.
	// (example value: "FF00FF" for magenta).
	Trans string `xml:"trans,attr"`

	//The image width in pixels.
	// (optional, used for tile index correction when the image changes).
	Width int32 `xml:"width,attr"`

	// The image height in pixels (optional).
	Height int32 `xml:"height,attr"`

	// Can contain: data (since 0.9.0)
	Data *Data `xml:"data"`
}

type Terrain struct {
	// The name of the terrain type.
	Name string `xml:"name,attr"`

	// The local tile-id of the tile that represents the terrain visually.
	Tile int32 `xml:"tile,attr"`

	// Can contain properties.
	Properties []Property `xml:"properties>property"`
}

type TilesetTile struct {
	// The local tile ID within its tileset.
	Id uint32 `xml:"id,attr"`

	// Defines the terrain type of each corner of the tile,
	// given as comma-separated indexes in the terrain types array
	// in the order top-left, top-right, bottom-left, bottom-right.
	// Leaving out a value means that corner has no terrain.
	// (optional) (since 0.9.0)
	Terrain string `xml:"terrain,attr"`

	// A percentage indicating the probability that this tile is
	// chosen when it competes with others while editing with
	// the terrain tool. (optional) (since 0.9.0)
	Probability float32 `xml:"probability,attr"`

	// Can contain properties.
	Properties []Property `xml:"properties>property"`

	// Can contain image (since 0.9.0).
	Image *Image `xml:"image"`
}

// All <tileset> tags shall occur before the first <layer> tag so that
// parsers may rely on having the tilesets before needing to resolve tiles.
type Layer struct {
	// The name of the layer.
	Name string `xml:"name,attr"`

	// The x coordinate of the layer in tiles. Defaults to 0 and
	// can no longer be changed in Tiled Qt.
	X int32 `xml:"x,attr"`

	// The y coordinate of the layer in tiles. Defaults to 0 and
	// can no longer be changed in Tiled Qt.
	Y int32 `xml:"y,attr"`

	// The width of the layer in tiles. Traditionally required, but
	// as of Tiled Qt always the same as the map width.
	Width int32 `xml:"width,attr"`

	// The height of the layer in tiles. Traditionally required, but
	// as of Tiled Qt always the same as the map height.
	Height int32 `xml:"height,attr"`

	// The opacity of the layer as a value from 0 to 1. Defaults to 1.
	Opacity float32 `xml:"opacity,attr"`

	// Whether the layer is shown (1) or hidden (0). Defaults to 1.
	Visible bool `xml:"visible,attr"`

	// Can contain properties.
	Properties []Property `xml:"properties>property"`

	// Can contain data.
	Data *Data `xml:"data"`
}

// When no encoding or compression is given, the tiles are stored as
// individual XML tile elements. Next to that, the easiest format
// to parse is the "csv" (comma separated values) format.
//
// The base64-encoded and optionally compressed layer data is somewhat
// more complicated to parse. First you need to base64-decode it, then you
// may need to decompress it. Now you have an array of bytes, which should
// be interpreted as an array of unsigned 32-bit integers using little-endian
// byte ordering.
//
// Whatever format you choose for your layer data, you will always end up
// with so called "global tile IDs" (gids). They are global, since they
// may refer to a tile from any of the tilesets used by the map. In order
// to find out from which tileset the tile is you need to find the tileset
// with the highest firstgid that is still lower or equal than the gid.
// The tilesets are always stored with increasing firstgids.
type Data struct {
	// The encoding used to encode the tile layer data.
	// When used, it can be "base64" and "csv" at the moment.
	Encoding string `xml:"encoding,attr"`

	// The compression used to compress the tile layer data.
	// Tiled Qt supports "gzip" and "zlib".
	Compression string `xml:"compression,attr"`

	// Can contain tile.
	RawTiles []DataTile `xml:"tile"`

	RawContents string `xml:",chardata"`
}

func (d *Data) Contents() string {
	return strings.TrimSpace(d.RawContents)
}

func (d *Data) base64Tiles() (tiles []DataTile, err error) {
	var (
		data  []byte
		buf   *bytes.Reader
		r     io.ReadCloser
		count int32
		gids  []uint32
	)
	if data, err = base64.StdEncoding.DecodeString(d.Contents()); err != nil {
		return
	}
	switch d.Compression {
	case "gzip":
		buf = bytes.NewReader(data)
		if r, err = gzip.NewReader(buf); err != nil {
			return
		}
		defer r.Close()
		if data, err = ioutil.ReadAll(r); err != nil {
			return
		}
	case "zlib":
		buf = bytes.NewReader(data)
		if r, err = zlib.NewReader(buf); err != nil {
			return
		}
		defer r.Close()
		if data, err = ioutil.ReadAll(r); err != nil {
			return
		}
	}
	buf = bytes.NewReader(data)
	count = int32(len(data) / binary.Size(count))
	gids = make([]uint32, count)
	if err = binary.Read(buf, binary.LittleEndian, &gids); err != nil {
		return
	}
	tiles = make([]DataTile, count)
	for i := 0; i < len(tiles); i++ {
		tiles[i].Gid = gids[i]
	}
	return
}

func (d *Data) csvTiles() (tiles []DataTile, err error) {
	err = fmt.Errorf("Not implemented")
	return
}

func (d *Data) Tiles() (tiles []DataTile, err error) {
	switch d.Encoding {
	case "base64":
		tiles, err = d.base64Tiles()
	case "csv":
		tiles, err = d.csvTiles()
	default:
		tiles = d.RawTiles
	}
	return
}

// Not to be confused with the tile element inside a tileset,
// this element defines the value of a single tile on a tile layer.
// This is however the most inefficient way of storing the tile layer data,
// and should generally be avoided.
type DataTile struct {
	// The global tile ID.
	Gid uint32 `xml:"gid,attr"`
}

// The object group is in fact a map layer,
// and is hence called "object layer" in Tiled Qt.
type ObjectGroup struct {
	// The name of the object group.
	Name string `xml:"name,attr"`

	// The color used to display the objects in this group.
	Color string `xml:"color,attr"`

	// The x coordinate of the object group in tiles.
	// Defaults to 0 and can no longer be changed in Tiled Qt.
	X int32 `xml:"x,attr"`

	// The y coordinate of the object group in tiles.
	// Defaults to 0 and can no longer be changed in Tiled Qt.
	Y int32 `xml:"y,attr"`

	// The width of the object group in tiles. Meaningless.
	Width int32 `xml:"width,attr"`

	// The height of the object group in tiles. Meaningless.
	Height int32 `xml:"height,attr"`

	// The opacity of the layer as a value from 0 to 1. Defaults to 1.
	Opacity float32 `xml:"opacity,attr"`

	// Whether the layer is shown (1) or hidden (0). Defaults to 1.
	Visible bool `xml:"visible,attr"`

	// Can contain properties.
	Properties []Property `xml:"properties>property"`

	// Can contain object.
	Objects []Object `xml:"object"`
}

// While tile layers are very suitable for anything repetitive
// aligned to the tile grid, sometimes you want to annotate
// your map with other information, not necessarily aligned to
// the grid. Hence the objects have their coordinates and size in
// pixels, but you can still easily align that to the grid when you want to.
//
// You generally use objects to add custom information to your
// tile map, such as spawn points, warps, exits, etc.
//
// When the object has a gid set, then it is represented by the
// image of the tile with that global ID. Currently that means width
// and height are ignored for such objects. The image alignment
// currently depends on the map orientation. In orthogonal orientation
// it's aligned to the bottom-left while in isometric it's aligned
// to the bottom-center.
type Object struct {
	// name: The name of the object. An arbitrary string.
	Name string `xml:"name,attr"`

	// type: The type of the object. An arbitrary string.
	Type string `xml:"type,attr"`

	// x: The x coordinate of the object in pixels.
	X int32 `xml:"x,attr"`

	// y: The y coordinate of the object in pixels.
	Y int32 `xml:"y,attr"`

	// width: The width of the object in pixels (defaults to 0).
	Width int32 `xml:"width,attr"`

	// height: The height of the object in pixels (defaults to 0).
	Height int32 `xml:"height,attr"`

	// rotation: The rotation of the object in degrees clockwise
	// (defaults to 0). (on git master)
	Rotation int32 `xml:"rotation,attr"`

	// gid: An reference to a tile (optional).
	Gid *uint32 `xml:"gid,attr"`

	// visible: Whether the object is shown (1) or hidden (0).
	// Defaults to 1. (since 0.9.0)
	Visible bool `xml:"visible,attr"`

	// Can contain properties.
	Properties []Property `xml:"properties>property"`

	// Can contain ellipse (since 0.9.0).
	Ellipse *Ellipse `xml:"ellipse"`

	// Can contain polygon.
	Polygon *Polygon `xml:"polygon"`

	// Can contain polyline.
	Polyline *Polyline `xml:"polyline"`

	// Can contain image.
	Image *Image `xml:"image"`
}

// Used to mark an object as an ellipse.
// The regular x, y, width, height attributes are used to
// determine the size of the ellipse.
type Ellipse struct{}

// Each polygon object is made up of a space-delimited list of x,y coordinates.
// The origin for these coordinates is the location of the parent object.
// By default, the first point is created as 0,0 denoting that the point
// will originate exactly where the object is placed.
type Polygon struct {
	RawPoints string `xml:"points,attr"`
}

// A polyline follows the same placement definition as a polygon object.
type Polyline struct {
	RawPoints string `xml:"points,attr"`
}

// A layer consisting of a single image.
type ImageLayer struct {
	// The name of the image layer.
	Name string `xml:"name,attr"`

	// The width of the image layer in tiles. Meaningless.
	Width int32 `xml:"width,attr"`

	// The height of the image layer in tiles. Meaningless.
	Height int32 `xml:"height,attr"`

	// opacity: The opacity of the layer as a value from 0 to 1.
	// Defaults to 1.
	Opacity float32 `xml:"opacity,attr"`

	// Whether the layer is shown (1) or hidden (0). Defaults to 1.
	Visible bool `xml:"visible,attr"`

	// Can contain properties.
	Properties []Property `xml:"properties>property"`

	// Can contain image.
	Image *Image `xml:"image"`
}

// When the property spans contains newlines, the current versions
// of Tiled Java and Tiled Qt will write out the value as characters
// contained inside the property element rather than as the value
// attribute. However, it is at the moment not really possible to
// edit properties consisting of multiple lines with Tiled.
//
// It is possible that a future version of the TMX format will switch
// to always saving property values inside the element rather than as
// an attribute.
type Property struct {
	// The name of the property.
	Name string `xml:"name,attr"`

	// The value of the property.
	Value string `xml:"value,attr"`
}

func ParseMapString(data string) (m *Map, err error) {
	m = &Map{}
	err = xml.Unmarshal([]byte(data), m)
	return
}
