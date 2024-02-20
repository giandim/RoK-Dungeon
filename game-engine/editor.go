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

var block [][]Tile

func renderGrid() {
	grid := fmt.Sprintf("<div class='grid col-%d'>", len(block[0]))

	for y := range block {
		for x, tile := range block[y] {
			secondLayer := ""
			cssClass := tile.layers[0].materialType + " " + tile.layers[0].direction

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
	isBlocking bool
}

func getButtons(this js.Value, args []js.Value) interface{} {
	buttons := [2]Button{
		{"wall", []string{"top", "left", "right", "tl", "tr"}, true},
		{"ground", []string{}, false},
	}

	resetButton := `<button onclick="createBlock()">Reset Block</button>`
	gridButton := `<button onclick="toggleGrid()">Toggle Grid</button>`
	var buttonGroup string

	for i, button := range buttons {
		buttonGroup += fmt.Sprintf(`<button onclick="showButtons(%d)">%s</button>`, i, button.buttonType)
		buttonGroup += fmt.Sprintf(`<div id="bgroup-%d" class="bgroup hidden">`, i)

		if len(button.directions) > 0 {
			for _, direction := range button.directions {
				buttonGroup += fmt.Sprintf(`<button class="%s %s" onclick="EDITOR.commandType='%s'; EDITOR.commandParams = {direction: '%s'};"></button>`, button.buttonType, direction, button.buttonType, direction)
			}
		} else {
			buttonGroup += fmt.Sprintf(`<button class="%s" onclick="EDITOR.commandType='%s'; EDITOR.commandParams = [];"></button>`, button.buttonType, button.buttonType)
		}
		buttonGroup += "</div>"
	}

	js.Global().Get("document").Call("getElementById", "command-panel").Set("innerHTML", resetButton+gridButton+buttonGroup)

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

func selectTile(this js.Value, args []js.Value) interface{} {
	fmt.Println(args)
	coordinates := strings.Split(args[0].String(), ",")

	y, _ := strconv.Atoi(coordinates[0])
	x, _ := strconv.Atoi(coordinates[1])

	block[y][x].layers[0].materialType = args[1].String()
	block[y][x].isBlocking = args[2].Bool()
	block[y][x].layers[0].direction = args[3].String()

	renderGrid()
	return nil
}

func registerCallbacks() {
	js.Global().Set("_EDITOR_getButtons", js.FuncOf(getButtons))
	js.Global().Set("_EDITOR_createBlock", js.FuncOf(createBlock))
	js.Global().Set("_EDITOR_selectTile", js.FuncOf(selectTile))
}

func main() {
	registerCallbacks()
	select {}
}
