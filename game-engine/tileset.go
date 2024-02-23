package main

import "fmt"

const (
	Wall      = "wall"
	Floor     = "floor"
	Collision = "collision"
	Item      = "item"
	Empty     = "empty"
)

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
		"item": {
			{id: 200, x: 0, y: 8},
			{id: 201, x: 1, y: 8},
			{id: 202, x: 2, y: 8},
			{id: 203, x: 3, y: 8},
			{id: 204, x: 4, y: 8},
			{id: 205, x: 5, y: 8},
			{id: 206, x: 6, y: 8},
			{id: 207, x: 7, y: 8},
			{id: 208, x: 8, y: 8},
			{id: 209, x: 9, y: 8},
			{id: 210, x: 0, y: 9},
			{id: 211, x: 1, y: 9},
			{id: 212, x: 7, y: 9},
			{id: 213, x: 8, y: 9},
			{id: 214, x: 7, y: 3},
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
