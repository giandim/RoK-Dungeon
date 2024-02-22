package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

type Layer struct {
	materialType string
	tileId       uint8
}

type Tile struct {
	isBlocking bool
	layers     [2]Layer
}

const (
	Wall      = "wall"
	Floor     = "floor"
	Collision = "collision"
	Item      = "item"
)

var block [][]Tile

func getTileTypes() map[string]bool {
	return map[string]bool{
		Collision: true,
		Wall:      true,
		Floor:     true,
		Item:      true,
	}
}

// Render the main grid for the editor
func renderGrid() {
	grid := "<div class='grid'>"

	for y := range block {
		for x, tile := range block[y] {
			secondLayer := ""
			cssClass := tile.layers[0].materialType

			var tileX, tileY int16

			if tile.isBlocking {
				cssClass += " collision"
			}

			// Add a second layer if present
			if tile.layers[1] != (Layer{}) {
				tileX, tileY, _ = getCoordinates(int16(tile.layers[1].tileId))
				secondLayer = fmt.Sprintf("<div onclick=\"selectTile('%d,%d')\" class=\"%s\" style=\"background-position: %dpx %dpx\"></div>", y, x, tile.layers[1].materialType, tileX*16*-5, tileY*16*-5)
			}

			// add a player placeholder
			if y == 8 && x == 4 {
				secondLayer = fmt.Sprintf("<div onclick=\"selectTile('%d,%d')\" class=\"player\"></div>", y, x)
			}

			tileX, tileY, _ = getCoordinates(int16(tile.layers[0].tileId))
			grid += fmt.Sprintf("<div onclick=\"selectTile('%d,%d')\" class=\"%s\" style=\"background-position: %dpx %dpx\">%s</div>", y, x, cssClass, tileX*16*-5, tileY*16*-5, secondLayer)
		}
	}

	grid += "</div>"

	js.Global().Get("document").Call("getElementById", "editor").Set("innerHTML", grid)
}

type Button struct {
	buttonType string
	directions []string
}

func getButtons(this js.Value, args []js.Value) interface{} {
	resetButton := `<button onclick="createBlock()">Reset Block</button>`
	gridButton := `<button onclick="toggleGrid()">Toggle Grid</button>`
	collisionButton := `<button onclick="EDITOR.tileType='collision'">Add collision</button>`
	layerButtons := `
    <div class="layer-buttons">
      <button onclick="selectLayer(0)" class="active">Layer 1</button>
      <button onclick="selectLayer(1)">Layer 2</button>
      <button onclick="selectLayer()">Collisions</button>
    </div>
  `

	var buttonGroup string

	for tileType, tiles := range getTileset() {
		buttonGroup += fmt.Sprintf(`<button onclick="showButtons('%s')">%s</button>`, tileType, tileType)
		buttonGroup += fmt.Sprintf(`<div id="bgroup-%s" class="bgroup hidden">`, tileType)

		for _, tile := range tiles {
			// TODO: remove that 50px and put it in a proper place
			buttonGroup += fmt.Sprintf(`<button class="tile" style="background-position: %dpx %dpx" onclick="EDITOR.tileType='%s'; EDITOR.commandParams = {tileId: %d};"></button>`, tile.x*-50, tile.y*-50, tileType, tile.id)
		}

		buttonGroup += "</div>"
	}

	js.Global().Get("document").Call("getElementById", "command-panel").Set("innerHTML", resetButton+gridButton+layerButtons+buttonGroup+collisionButton)

	return nil
}

func createBlock(this js.Value, args []js.Value) interface{} {
	block = make([][]Tile, 9)
	defaultTile := Tile{layers: [2]Layer{{materialType: "floor", tileId: 150}, {}}, isBlocking: false}

	for i := range block {
		block[i] = make([]Tile, 9)

		for j := range block[i] {
			block[i][j] = defaultTile
		}
	}

	renderGrid()
	return nil
}

// args[0] -> coordinates (0,0)
// args[1] -> tileType (wall, ground, collision)
// args[2] -> tileId
// args[3] -> layer (0,1)
func setTile(this js.Value, args []js.Value) interface{} {
	fmt.Println(args)

	coordinates := strings.Split(args[0].String(), ",")
	tileType := args[1].String()
	layer := 0
	tileId := uint8(0)

	if !args[2].IsUndefined() {
		tileId = uint8(args[2].Int())
	}

	if !args[3].IsUndefined() {
		layer = args[3].Int()
	}

	y, _ := strconv.Atoi(coordinates[0])
	x, _ := strconv.Atoi(coordinates[1])

	if _, ok := getTileTypes()[tileType]; !ok {
		return nil
	}

	if tileType != Collision {
		block[y][x].layers[layer].materialType = tileType
		block[y][x].layers[layer].tileId = tileId
	} else {
		block[y][x].isBlocking = !block[y][x].isBlocking
	}

	renderGrid()
	return nil
}

func registerCallbacks() {
	js.Global().Set("_EDITOR_getButtons", js.FuncOf(getButtons))
	js.Global().Set("_EDITOR_createBlock", js.FuncOf(createBlock))
	js.Global().Set("_EDITOR_setTile", js.FuncOf(setTile))
}

func main() {
	registerCallbacks()
	select {}
}
