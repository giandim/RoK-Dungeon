package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

type Tile struct {
	tileType   string
	direction  string
	isBlocking bool
}

var block [][]Tile

func renderGrid() {
	grid := fmt.Sprintf("<div class='grid col-%d'>", len(block[0]))

	for y := range block {
		for x, tile := range block[y] {
			cssClass := tile.tileType + " " + tile.direction

			// add a player placeholder
			if y == 8 && x == 4 {
				cssClass = "player"
			}

			grid += fmt.Sprintf("<div onclick=\"selectTile('%d,%d')\" class=\"%s\"></div>", y, x, cssClass)
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

	jsButtons := js.Global().Get("Array").New(len(buttons))

	for i, button := range buttons {
		jsButton := js.Global().Get("Object").New()
		jsButton.Set("type", button.buttonType)
		jsButton.Set("directions", jsSliceOf(button.directions))
		jsButton.Set("blocking", button.isBlocking)
		jsButtons.SetIndex(i, jsButton)
	}

	return jsButtons
}

func createBlock(this js.Value, args []js.Value) interface{} {
	block = make([][]Tile, 9)
	defaultTile := Tile{tileType: "ground", direction: "", isBlocking: false}

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

	block[y][x].tileType = args[1].String()
	block[y][x].isBlocking = args[2].Bool()
	block[y][x].direction = args[3].String()

	renderGrid()

	// retrieves the JavaScript array constructor and create a new one
	jsTile := js.Global().Get("Array").New(len(block))

	for i, row := range block {
		jsRow := js.Global().Get("Array").New(len(row))
		for j, tile := range row {
			// set value at the index
			jsRow.SetIndex(j, tile.isBlocking)
		}
		jsTile.SetIndex(i, jsRow)
	}

	return jsTile
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
