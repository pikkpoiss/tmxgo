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

func TestParseMapString(t *testing.T) {
	var (
		m   *Map
		err error
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
}
