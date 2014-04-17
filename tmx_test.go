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
	"fmt"
	"testing"
)

const TEST_MAP = `
<?xml version="1.0" encoding="UTF-8"?>
<map version="1.0" orientation="orthogonal" width="71" height="40" tilewidth="16" tileheight="16">
 <properties>
  <property name="time1" value="16"/>
  <property name="time2" value="9"/>
  <property name="time3" value="6"/>
 </properties>
 <tileset firstgid="1" name="sprites32" tilewidth="32" tileheight="32">
  <image source="../textures/sprites32.png" width="512" height="64"/>
 </tileset>
 <tileset firstgid="33" name="sprites16" tilewidth="16" tileheight="16">
  <image source="../textures/sprites16.png" width="256" height="32"/>
 </tileset>
 <tileset firstgid="65" name="stars" tilewidth="16" tileheight="16">
  <image source="../textures/stars.png" width="64" height="16"/>
 </tileset>
 <layer name="Tile Layer 3" width="71" height="40">
  <data encoding="base64" compression="zlib">
   eJzt2MsKwjAQheEKvog7r4u+/8uZhQsrTh1wwtz+DwJdhJKeJs00ywL0dR7t5j2IoO6jPbwHEVSVbKzm//p2XSWb62gX43tWyWaGGXmjt/V3l5SO3gNASFbz3asuPij7ea5r9nAZ2WxVrItnIBsZdfFcp1f7V8WzJKtsPudwlvp1b5xW2VRUNZuKa9wKe7iMbLaoi3XIRpalLv52VsE+IsvyXhFbln+MCjp+z7TP3HEP137DO2azh7pYh2xk1E+I4glYsQ1i
  </data>
 </layer>
 <layer name="Stars" width="71" height="40">
  <data encoding="base64" compression="zlib">
   eJztl1sOhTAIRN2auv89Gf9MbC0FLK85n1dvAjNTsNtGYyf+lpXDugBjLPof5av3nFNrdX+jIfELXv8HZyec6lWsI3KWqLpn2vOR/brheMHtecW5bNWWKW/WQEuQGeS7Hqt3+HMPIm+gCpHvZeBNa3ZFvw9VADtHj56WvVmHGegPynnAXPMPPAKrQeb84fn7Rrs2z71qk7FXSU8zsyejdtZQ9R+9x/XGytNsWaqyw73cPb/0pmYrumfZztAIrq8zOkkzIfm/l7MVGYqGFwJPD+s=
  </data>
 </layer>
</map>
`

const TEST_TILES_FROM_LAYER_MAP = `
<?xml version="1.0" encoding="UTF-8"?>
<map version="1.0" orientation="orthogonal" width="2" height="2" tilewidth="16" tileheight="16">
 <tileset firstgid="1" name="sprites1" tilewidth="16" tileheight="16">
  <image source="../textures/sprites1.png" width="64" height="16"/>
 </tileset>
 <tileset firstgid="5" name="sprites2" tilewidth="16" tileheight="16">
  <image source="../textures/sprites2.png" width="64" height="16"/>
 </tileset>
 <layer name="layer1" width="2" height="2">
  <data>
   <tile gid="1" />
   <tile gid="5" />
   <tile gid="2" />
   <tile gid="6" />
  </data>
 </layer>
 <layer name="layer2" width="2" height="2">
  <data>
   <tile gid="2147483649" />
   <tile gid="1073741827" />
   <tile gid="536870916" />
   <tile gid="2684354574" />
  </data>
 </layer>
</map>
`

func TestParseGid(t *testing.T) {
	var (
		val uint32
		fh  bool
		fv  bool
		fd  bool
		id  uint32
	)
	type testcase struct {
		Input string
		Id    uint32
		Fh    bool
		Fv    bool
		Fd    bool
	}
	tests := []testcase{
		testcase{"10000000000000000000000000000001", 1, true, false, false},
		testcase{"01000000000000000000000000000011", 3, false, true, false},
		testcase{"00100000000000000000000000000100", 4, false, false, true},
		testcase{"10100000000000000000000000001110", 14, true, false, true},
	}
	for i := 0; i < len(tests); i++ {
		c := tests[i]
		if _, err := fmt.Sscanf(c.Input, "%b", &val); err != nil {
			t.Fatalf("Invalid Gid: %v", err)
		}
		id, fh, fv, fd = parseGid(val)
		if id != c.Id || fh != c.Fh || fv != c.Fv || fd != c.Fd {
			t.Errorf("Gid parsed wrong: %v %v %v %v %v", id, fh, fv, fd, c)
		}
	}
}

func TestTilesFromLayer(t *testing.T) {
	var (
		m     *Map
		tiles []Tile
		err   error
	)
	if m, err = ParseMapString(TEST_TILES_FROM_LAYER_MAP); err != nil {
		t.Errorf("Could not parse: %v", err)
	}
	if tiles, err = m.TilesFromLayerIndex(0); err != nil {
		t.Fatalf("Could not get layer 0")
	}
	if len(tiles) != 4 {
		t.Fatalf("Did not have enough tiles")
	}
	if tiles[0].Index != 0 {
		t.Errorf("Wrong index: %v", tiles[0].Index)
	}
	if tiles[0].FlipHorz == true {
		t.Errorf("FlipHorz parsed incorrectly")
	}
	if tiles[1].Index != 0 {
		t.Errorf("Wrong index: %v", tiles[1].Index)
	}
	if tiles[1].FlipVert == true {
		t.Errorf("FlipVert parsed incorrectly")
	}
	if tiles[2].Index != 1 {
		t.Errorf("Wrong index: %v", tiles[2].Index)
	}
	if tiles[2].FlipDiag == true {
		t.Errorf("FlipDiag parsed incorrectly")
	}
	if tiles[3].Index != 1 {
		t.Errorf("Wrong index: %v", tiles[3].Index)
	}
	if tiles[3].FlipHorz == true || tiles[3].FlipDiag == true {
		t.Errorf("FlipHorz & FlipDiag parsed incorrectly")
	}
	if tiles, err = m.TilesFromLayerName("layer2"); err != nil {
		t.Fatalf("Could not get layer 'layer2'")
	}
	if len(tiles) != 4 {
		t.Fatalf("Did not have enough tiles")
	}
	if tiles[0].Index != 1-1 {
		t.Errorf("Wrong index: %v", tiles[0].Index)
	}
	if tiles[0].FlipHorz == false {
		t.Errorf("FlipHorz parsed incorrectly")
	}
	if tiles[1].Index != 3-1 {
		t.Errorf("Wrong index: %v", tiles[1].Index)
	}
	if tiles[1].FlipVert == false {
		t.Errorf("FlipVert parsed incorrectly")
	}
	if tiles[2].Index != 4-1 {
		t.Errorf("Wrong index: %v", tiles[2].Index)
	}
	if tiles[2].FlipDiag == false {
		t.Errorf("FlipDiag parsed incorrectly")
	}
	if tiles[3].Index != 14-5 {
		t.Errorf("Wrong index: %v", tiles[3].Index)
	}
	if tiles[3].FlipHorz == false || tiles[3].FlipDiag == false {
		t.Errorf("FlipHorz & FlipDiag parsed incorrectly")
	}
	if tiles[0].Tileset.Name != "sprites1" {
		t.Errorf("Invalid tileset: %v", tiles[0].Tileset.Name)
	}
	if tiles[1].Tileset.Name != "sprites1" {
		t.Errorf("Invalid tileset: %v", tiles[1].Tileset.Name)
	}
	if tiles[2].Tileset.Name != "sprites1" {
		t.Errorf("Invalid tileset: %v", tiles[2].Tileset.Name)
	}
	if tiles[3].Tileset.Name != "sprites2" {
		t.Errorf("Invalid tileset: %v", tiles[3].Tileset.Name)
	}
}

func TestParseMapString(t *testing.T) {
	var (
		m     *Map
		tiles []DataTile
		err   error
	)
	if m, err = ParseMapString(TEST_MAP); err != nil {
		t.Errorf("Could not parse: %v", err)
	}
	if m.Version != "1.0" {
		t.Errorf("Invalid version: %v", m.Version)
	}
	if m.Orientation != "orthogonal" {
		t.Errorf("Invalid orientation: %v", m.Orientation)
	}
	if m.Width != 71 {
		t.Errorf("Invalid width: %v", m.Width)
	}
	if m.Height != 40 {
		t.Errorf("Invalid height: %v", m.Height)
	}
	if m.TileWidth != 16 {
		t.Errorf("Invalid tilewidth: %v", m.TileWidth)
	}
	if m.TileHeight != 16 {
		t.Errorf("Invalid tileheight: %v", m.TileHeight)
	}
	if len(m.Properties) != 3 {
		t.Fatalf("Not enough properties: %v", len(m.Properties))
	}
	if m.Properties[0].Name != "time1" {
		t.Errorf("Invalid property name: %v", m.Properties[0].Name)
	}
	if m.Properties[0].Value != "16" {
		t.Errorf("Invalid property value: %v", m.Properties[0].Value)
	}
	if len(m.Tilesets) != 3 {
		t.Fatalf("Not enough tilesets: %v", len(m.Tilesets))
	}
	if m.Tilesets[0].FirstGid != 1 {
		t.Errorf("Invalid firstgid: %v", m.Tilesets[0].FirstGid)
	}
	if m.Tilesets[0].Name != "sprites32" {
		t.Errorf("Invalid name: %v", m.Tilesets[0].Name)
	}
	if m.Tilesets[0].TileWidth != 32 {
		t.Errorf("Invalid tilewidth: %v", m.Tilesets[0].TileWidth)
	}
	if m.Tilesets[0].TileHeight != 32 {
		t.Errorf("Invalid tileheight: %v", m.Tilesets[0].TileHeight)
	}
	if m.Tilesets[0].Image == nil {
		t.Fatalf("No image")
	}
	if m.Tilesets[0].Image.Source != "../textures/sprites32.png" {
		t.Errorf("Invalid source: %v", m.Tilesets[0].Image.Source)
	}
	if m.Tilesets[0].Image.Width != 512 {
		t.Errorf("Invalid width: %v", m.Tilesets[0].Image.Width)
	}
	if m.Tilesets[0].Image.Height != 64 {
		t.Errorf("Invalid height: %v", m.Tilesets[0].Image.Height)
	}
	if len(m.Layers) != 2 {
		t.Fatalf("Not enough layers: %v", len(m.Layers))
	}
	if m.Layers[0].Name != "Tile Layer 3" {
		t.Errorf("Invalid name: %v", m.Layers[0].Name)
	}
	if m.Layers[0].Width != 71 {
		t.Errorf("Invalid width: %v", m.Layers[0].Width)
	}
	if m.Layers[0].Height != 40 {
		t.Errorf("Invalid height: %v", m.Layers[0].Height)
	}
	if m.Layers[0].Data == nil {
		t.Fatalf("No data")
	}
	if m.Layers[0].Data.Encoding != "base64" {
		t.Errorf("Invalid encoding: %v", m.Layers[0].Data.Encoding)
	}
	if m.Layers[0].Data.Compression != "zlib" {
		t.Errorf("Invalid compression: %v", m.Layers[0].Data.Compression)
	}
	if m.Layers[0].Data.Contents()[0:10] != "eJzt2MsKwj" {
		t.Errorf("Invalid data string: %v", m.Layers[0].Data.Contents()[0:10])
	}
	if tiles, err = m.Layers[1].Data.Tiles(); err != nil {
		t.Fatalf("Invalid tiles: %v", err)
	}
	if len(tiles) != 2840 {
		t.Errorf("Invalid tiles length: %v", len(tiles))
	}
	if tiles[10].Gid != 65 {
		t.Errorf("Invalid tile gid: %v", tiles[9].Gid)
	}
}
