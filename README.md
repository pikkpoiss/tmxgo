tmxgo
=====

TMX map parser for Go.  Can parse the map files created by the [Tiled map editor](http://www.mapeditor.org/)

Based off of the TMX specification from <https://github.com/bjorn/tiled/wiki/TMX-Map-Format>.

Supports:

  * Gzip compression
  * Zlib compression
  * Base64 encoded tiles
  * Unencoded tile elements
  * Serializing a map back to a string (for edit + save)

TODO:

  * Support CSV encoding.
  * Unit tests for full spec.

## Documentation

<https://godoc.org/github.com/kurrik/tmxgo>
