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
	"encoding/xml"
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

	// Can contain imagelayer.
}

type Tileset struct {
	// The first global tile ID of this tileset.
	// (this global ID maps to the first tile in this tileset).
	FirstGid int32 `xml:"firstgid,attr"`

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
	Id int32 `xml:"id,attr"`

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
	Tiles []Tile `xml:"tile"`

	RawContents string `xml:",chardata"`
}

func (d *Data) Contents() string {
	return strings.TrimSpace(d.RawContents)
}

// Not to be confused with the tile element inside a tileset,
// this element defines the value of a single tile on a tile layer.
// This is however the most inefficient way of storing the tile layer data,
// and should generally be avoided.
type Tile struct {
	// The global tile ID.
	Gid int32 `xml:"gid,attr"`
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
