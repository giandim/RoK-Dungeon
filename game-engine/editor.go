package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	fetch "marwan.io/wasm-fetch"
)

type Layer struct {
	MaterialType string `json:"materialType"`
	TileId       uint8  `json:"tileId"`
}

type Tile struct {
	IsBlocking bool     `json:"isBlocking"`
	Layers     [2]Layer `json:"layers"`
}

type Block struct {
	id          string
	Tiles       [][]Tile `json:"tiles"`
	Connections [4]bool  `json:"connections"`
}

const (
	DefaultScale              = 5
	DefaultTileSize           = 16
	TileCoordinatesMultiplier = DefaultTileSize * -DefaultScale
	ButtonTileMultiplier      = -50
)

var block Block

// Render the main grid for the editor
func renderGrid() {
	grid := `<div class="grid">`
	tiles := &block.Tiles

	for y, row := range *tiles {
		for x, tile := range row {
			cssClass := tile.Layers[0].MaterialType

			var tileX, tileY int16

			if tile.IsBlocking {
				cssClass += " collision"
			}

			var secondLayer string
			// Add a second layer if present
			if tile.Layers[1] != (Layer{}) {
				tileX, tileY, _ = getCoordinates(int16(tile.Layers[1].TileId))
				secondLayer = fmt.Sprintf(`<div class="%s" style="background-position: %dpx %dpx"></div>`, tile.Layers[1].MaterialType, tileX*TileCoordinatesMultiplier, tileY*TileCoordinatesMultiplier)
			}

			// add a player placeholder
			if y == 8 && x == 4 {
				secondLayer = fmt.Sprintf(`<div onclick="selectTile('%d,%d')" class="player"></div>`, y, x)
			}

			tileX, tileY, _ = getCoordinates(int16(tile.Layers[0].TileId))

			// Append the tile HTML to the grid
			grid += fmt.Sprintf(`<div onclick="selectTile('%d,%d')" class="%s" style="background-position: %dpx %dpx">%s</div>`, y, x, cssClass, tileX*TileCoordinatesMultiplier, tileY*TileCoordinatesMultiplier, secondLayer)
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
	collisionButton := fmt.Sprintf(`<button onclick="EDITOR.tileType='%s'">Add collision</button>`, Collision)
	removeButton := fmt.Sprintf(`<button onclick="EDITOR.tileType='%s'">Remove tile</button>`, Empty)
	saveButton := `<button onclick="saveBlock()">Save Block</button>`

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
			buttonGroup += fmt.Sprintf(`<button class="tile" style="background-position: %dpx %dpx" onclick="EDITOR.tileType='%s'; EDITOR.commandParams = {tileId: %d};"></button>`, tile.x*ButtonTileMultiplier, tile.y*ButtonTileMultiplier, tileType, tile.id)
		}

		buttonGroup += "</div>"
	}

	connections := `
    <div class="connections">
      Tile Connections<br>
      <input type="checkbox" id="up" name="connections" value="up"><label for="up">top</label> 
      <input type="checkbox" id="right" name="connections" value="right"><label for="right">right</label>
      <input type="checkbox" id="down" name="connections" value="down"><label for="down">bottom</label>
      <input type="checkbox" id="left" name="connections" value="left"><label for="left">left</label>
    </div>
  `

	js.Global().Get("document").Call("getElementById", "command-panel").Set("innerHTML", saveButton+resetButton+gridButton+layerButtons+buttonGroup+collisionButton+removeButton+connections)

	return nil
}

func createBlock(this js.Value, args []js.Value) interface{} {
	tiles := &block.Tiles
	*tiles = make([][]Tile, 9)
	defaultTile := Tile{Layers: [2]Layer{{MaterialType: "floor", TileId: 150}, {}}, IsBlocking: false}

	for i := range *tiles {
		(*tiles)[i] = make([]Tile, 9)

		for j := range (*tiles)[i] {
			(*tiles)[i][j] = defaultTile
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
	coordinates := strings.Split(args[0].String(), ",")
	tileType := args[1].String()
	layer := 0
	tileId := uint8(0)
	tiles := &block.Tiles

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

	if tileType == Empty {
		// Clear the layer if the tile type is Empty
		(*tiles)[y][x].Layers[layer] = Layer{}
	} else if tileType != Collision && !args[3].IsUndefined() {
		(*tiles)[y][x].Layers[layer].MaterialType = tileType
		(*tiles)[y][x].Layers[layer].TileId = tileId
	} else if tileType == Collision {
		(*tiles)[y][x].IsBlocking = !(*tiles)[y][x].IsBlocking
	} else {
		return nil
	}

	renderGrid()
	return nil
}

// args[0] -> connections
func saveBlock(this js.Value, args []js.Value) interface{} {
	go func() {
		for i := 0; i < 4; i++ {
			block.Connections[i] = args[0].Index(i).Bool()
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		data, err := json.Marshal(&block)
		if err != nil {
			// Handle the error, such as logging or returning an error to the caller
			fmt.Printf("Error marshaling block data: %v", err)
			return
		}
		resp, _ := fetch.Fetch("http://localhost:8081/api/blocks", &fetch.Opts{
			Body:   bytes.NewBuffer(data),
			Method: fetch.MethodPost,
			Signal: ctx,
		})

		block.id = string(resp.Body)
		js.Global().Get("alert").Invoke("Block Saved!")
	}()

	return nil
}

func getBlocks(this js.Value, args []js.Value) interface{} {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		defer cancel()

		resp, err := fetch.Fetch("http://localhost:8081/api/blocks", &fetch.Opts{
			Method: fetch.MethodGet,
			Signal: ctx,
		})

		var blocks []string

		if err = json.Unmarshal(resp.Body, &blocks); err != nil {
			fmt.Println("Error parsing blocks json: ", err)
			return
		}

		blocksHtml := ""
		for _, block := range blocks {
			blocksHtml += fmt.Sprintf(`<button onclick="_EDITOR_loadBlock('%s');_EDITOR_getButtons()">%s</button>`, block, block)
		}

		js.Global().Get("document").Call("getElementById", "command-panel").Set("innerHTML", blocksHtml)
	}()
	return nil
}

// args[0] -> block id
func loadBlock(this js.Value, args []js.Value) interface{} {
	blockId := args[0].String()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		defer cancel()

		resp, err := fetch.Fetch("http://localhost:8081/api/blocks/"+blockId, &fetch.Opts{
			Method: fetch.MethodGet,
			Signal: ctx,
		})

		if err = json.Unmarshal(resp.Body, &block); err != nil {
			fmt.Println("Error Unmarshaling block: ", err)
			return
		}
		renderGrid()
	}()
	return nil
}

func registerCallbacks() {
	js.Global().Set("_EDITOR_getButtons", js.FuncOf(getButtons))
	js.Global().Set("_EDITOR_createBlock", js.FuncOf(createBlock))
	js.Global().Set("_EDITOR_saveBlock", js.FuncOf(saveBlock))
	js.Global().Set("_EDITOR_getBlocks", js.FuncOf(getBlocks))
	js.Global().Set("_EDITOR_loadBlock", js.FuncOf(loadBlock))
	js.Global().Set("_EDITOR_setTile", js.FuncOf(setTile))
}

func main() {
	registerCallbacks()
	select {}
}
