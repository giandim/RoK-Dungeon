package main

import "fmt"

type Tileset struct {
	id int16
	x  int16
	y  int16
}

func getTileset() map[string][]Tileset {
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
		},
		"floor": {
			{id: 150, x: 9, y: 7},
			{id: 151, x: 6, y: 0},
			{id: 152, x: 7, y: 0},
			{id: 153, x: 8, y: 0},
			{id: 154, x: 9, y: 0},
		},
	}
}

func getCoordinates(tilesetId int16) (int16, int16, error) {
	for _, tilesets := range getTileset() {
		for _, tileset := range tilesets {
			if tileset.id == tilesetId {
				return tileset.x, tileset.y, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("Tileset with id %d not found", tilesetId)
}
