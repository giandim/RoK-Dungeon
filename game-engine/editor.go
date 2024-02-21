package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

type Layer struct {
	materialType string
	direction    string
}

type Tile struct {
	isBlocking bool
	layers     [2]Layer
}

const (
	Wall      = "wall"
	Ground    = "ground"
	Collision = "collision"
)

var block [][]Tile

func getTileTypes() map[string]bool {
	return map[string]bool{
		Collision: true,
		Wall:      true,
		Ground:    true,
	}
}

// Render the main grid for the editor
func renderGrid() {
	grid := fmt.Sprintf("<div class='grid col-%d'>", len(block[0]))

	for y := range block {
		for x, tile := range block[y] {
			secondLayer := ""
			cssClass := tile.layers[0].materialType + " " + tile.layers[0].direction

			if tile.isBlocking {
				cssClass += " collision"
			}

			// Add a second layer if present
			if tile.layers[1] != (Layer{}) {
				secondLayer = fmt.Sprintf("<div onclick=\"selectTile('%d,%d')\" class=\"%s\"></div>", y, x, tile.layers[1].materialType)
			}

			// add a player placeholder
			if y == 8 && x == 4 {
				cssClass = "player"
			}

			grid += fmt.Sprintf("<div onclick=\"selectTile('%d,%d')\" class=\"%s\">%s</div>", y, x, cssClass, secondLayer)
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
	buttons := [2]Button{
		{"wall", []string{"top", "left", "right", "tl", "tr"}},
		{"ground", []string{}},
	}

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

	for i, button := range buttons {
		buttonGroup += fmt.Sprintf(`<button onclick="showButtons(%d)">%s</button>`, i, button.buttonType)
		buttonGroup += fmt.Sprintf(`<div id="bgroup-%d" class="bgroup hidden">`, i)

		if len(button.directions) > 0 {
			for _, direction := range button.directions {
				buttonGroup += fmt.Sprintf(`<button class="%s %s" onclick="EDITOR.tileType='%s'; EDITOR.commandParams = {direction: '%s'};"></button>`, button.buttonType, direction, button.buttonType, direction)
			}
		} else {
			buttonGroup += fmt.Sprintf(`<button class="%s" onclick="EDITOR.tileType='%s'; EDITOR.commandParams = [];"></button>`, button.buttonType, button.buttonType)
		}
		buttonGroup += "</div>"
	}

	js.Global().Get("document").Call("getElementById", "command-panel").Set("innerHTML", resetButton+gridButton+layerButtons+buttonGroup+collisionButton)

	return nil
}

func createBlock(this js.Value, args []js.Value) interface{} {
	block = make([][]Tile, 9)
	defaultTile := Tile{layers: [2]Layer{{materialType: "ground", direction: ""}, {}}, isBlocking: false}

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
// args[2] -> direction (top, left, right, tl, tr)
// args[3] -> layer (0,1)
func setTile(this js.Value, args []js.Value) interface{} {
	fmt.Println(args)

	coordinates := strings.Split(args[0].String(), ",")
	tileType := args[1].String()
	direction := args[2].String()
	layer := 0

	if !args[3].IsUndefined() {
		layer = args[3].Int()
	}

	y, _ := strconv.Atoi(coordinates[0])
	x, _ := strconv.Atoi(coordinates[1])

	_, ok := getTileTypes()[tileType]

	if !ok {
		return nil
	}

	if tileType != Collision {
		block[y][x].layers[layer].materialType = tileType
		block[y][x].layers[layer].direction = direction
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
