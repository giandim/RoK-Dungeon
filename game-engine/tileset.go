package main

import "fmt"

const (
	Wall      = "wall"
	Floor     = "floor"
	Collision = "collision"
	Item      = "item"
	Empty     = "empty"
)

var tilesets = [5]string{"dungeon", "castle", "town", "forest", "world"}

func getTileTypes() map[string]bool {
	return map[string]bool{
		Collision: true,
		Wall:      true,
		Floor:     true,
		Item:      true,
		Empty:     true,
	}
}

type Tileset struct {
	id           int16
	x            int16
	y            int16
	hasAnimation bool
}

func getTileset(tilesetId int) map[string][]Tileset {
	if tilesetId > len(tilesets) {
		return nil
	}

	switch tilesets[tilesetId] {
	case "dungeon":
		return getDungeonTileset()
	case "forest":
		return getForestTileset()
	default:
		return nil
	}
}

func getDungeonTileset() map[string][]Tileset {
	return map[string][]Tileset{
		"wall": {
			{id: 100, x: 0, y: 0},
			{id: 101, x: 0, y: 1},
			{id: 102, x: 0, y: 4},
			{id: 103, x: 1, y: 0},
			{id: 104, x: 1, y: 4},
			{id: 105, x: 5, y: 0},
			{id: 106, x: 5, y: 1},
			{id: 107, x: 5, y: 4},
			{id: 108, x: 0, y: 5},
			{id: 109, x: 5, y: 5},
		},
		"floor": {
			{id: 200, x: 9, y: 7},
			{id: 201, x: 6, y: 0},
			{id: 202, x: 7, y: 0},
			{id: 203, x: 8, y: 0},
			{id: 204, x: 9, y: 0},
			{id: 205, x: 8, y: 7},
		},
		"item": {
			{id: 300, x: 0, y: 8},
			{id: 301, x: 1, y: 8},
			{id: 302, x: 2, y: 8},
			{id: 303, x: 3, y: 8},
			{id: 304, x: 4, y: 8},
			{id: 305, x: 5, y: 8},
			{id: 306, x: 6, y: 8},
			{id: 307, x: 7, y: 8},
			{id: 308, x: 8, y: 8},
			{id: 309, x: 9, y: 8},
			{id: 310, x: 0, y: 9},
			{id: 311, x: 1, y: 9},
			{id: 312, x: 7, y: 9},
			{id: 313, x: 8, y: 9},
			{id: 314, x: 7, y: 3},
			{id: 315, x: 7, y: 7},
		},
	}
}

func getForestTileset() map[string][]Tileset {
	return map[string][]Tileset{
		"floor": {
			{id: 200, x: 7, y: 0},
			{id: 201, x: 8, y: 0},
			{id: 202, x: 7, y: 1},
			{id: 203, x: 8, y: 1},
			{id: 204, x: 1, y: 0},
			{id: 205, x: 0, y: 1},
			{id: 206, x: 2, y: 1},
			{id: 207, x: 1, y: 2},
			{id: 208, x: 3, y: 0},
			{id: 209, x: 4, y: 0},
			{id: 210, x: 3, y: 1},
			{id: 211, x: 4, y: 1},
			{id: 212, x: 5, y: 0},
			{id: 213, x: 6, y: 0},
			{id: 214, x: 5, y: 1},
			{id: 215, x: 6, y: 1},
		},
		"item": {
			{id: 300, x: 10, y: 0},
			{id: 301, x: 11, y: 0},
			{id: 302, x: 12, y: 0},
			{id: 303, x: 10, y: 1},
			{id: 304, x: 12, y: 1},
			{id: 305, x: 10, y: 2},
			{id: 306, x: 11, y: 2},
			{id: 307, x: 12, y: 2},
		},
	}
}

func getCoordinates(tilesetId int, tileId int16) (int16, int16, error) {
	for _, tiles := range getTileset(tilesetId) {
		for _, tile := range tiles {
			if tile.id == tileId {
				return tile.x, tile.y, nil
			}
		}
	}
	return -1, -1, fmt.Errorf("Tileset with id %d not found", tilesetId)
}
