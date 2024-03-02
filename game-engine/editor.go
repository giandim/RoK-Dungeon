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

const (
	DefaultScale         = 5
	DefaultTileSize      = 16
	ButtonTileMultiplier = -50
	DefaultGridSize      = 9
)

type Block struct {
	id          string
	Tiles       [][]Tile `json:"tiles"`
	Connections [4]bool  `json:"connections"`
}

type Tile struct {
	IsBlocking bool    `json:"isBlocking"`
	Layers     []Layer `json:"layers"`
}

type Layer struct {
	MaterialType string `json:"materialType"`
	TileId       uint8  `json:"tileId"`
}

var (
	block        Block
	currentScale int16 = DefaultScale
)

func main() {
	registerCallbacks()
	select {}
}

// TODO: add a duplicate block
// TODO: add a update block
// TODO: add a remove block
func registerCallbacks() {
	js.Global().Set("_EDITOR_getButtons", js.FuncOf(getButtons))
	js.Global().Set("_EDITOR_createBlock", js.FuncOf(createBlock))
	js.Global().Set("_EDITOR_saveBlock", js.FuncOf(saveBlock))
	js.Global().Set("_EDITOR_getBlocks", js.FuncOf(getBlocks))
	js.Global().Set("_EDITOR_loadBlock", js.FuncOf(loadBlock))
	js.Global().Set("_EDITOR_setTile", js.FuncOf(setTile))
	js.Global().Set("_EDITOR_changeZoom", js.FuncOf(changeZoom))
	js.Global().Set("_EDITOR_addLayer", js.FuncOf(addLayer))
	js.Global().Set("_EDITOR_renderLayersSection", js.FuncOf(renderLayersSection))
}

// Render the main grid for the editor
// grids[0] -> collisions and grid
// grids[1] -> layer 1 always rendered
func renderGrid() {
	grids := make([]string, len(block.Tiles[0][0].Layers)+1)
	tiles := &block.Tiles

	for y, row := range *tiles {
		// Render each tile
		for x, tile := range row {
			cssClass := ""

			var tileX, tileY int16
			tileX, tileY, _ = getCoordinates(int16(tile.Layers[0].TileId))

			// Render the grid and collision tile
			if tile.IsBlocking {
				cssClass = "collision"
			}

			grids[0] += fmt.Sprintf(`<div onclick="setTile('%d,%d')" class="%s"></div>`, y, x, cssClass)

			// The first layer is always rendered
			cssClass = tile.Layers[0].MaterialType

			// add a player placeholder
			if y == 4 && x == 4 {
				// secondLayer = fmt.Sprintf(`<div onclick="setTile('%d,%d')" class="player"></div>`, y, x)
			}

			// Append the tile HTML to the grid
			grids[1] += fmt.Sprintf(`<div class="%s" style="background-position: %dpx %dpx"></div>`, cssClass, tileX*currentScale*-DefaultTileSize, tileY*currentScale*-DefaultTileSize)

			// The other layers are only rendered if they exist
			for i := 1; i < len(tile.Layers); i++ {
				fmt.Println("layer i:", i)
				fmt.Println("layer:", tile.Layers[i])

				layer := tile.Layers[i]

				if layer.TileId > 0 {
					tileX, tileY, _ = getCoordinates(int16(tile.Layers[i].TileId))
					tileEl := fmt.Sprintf(`<div class="%s" style="background-position: %dpx %dpx; top: %dpx; left: %dpx"></div>`,
						tile.Layers[i].MaterialType, tileX*currentScale*-DefaultTileSize, tileY*currentScale*-DefaultTileSize, int16(y)*currentScale*DefaultTileSize, int16(x)*currentScale*DefaultTileSize)

					grids[i+1] += tileEl
				}
			}

		}
	}

	js.Global().Get("document").Call("getElementById", "editor").Set("innerHTML", "")

	// Render the grids
	for i, grid := range grids {
		gridEl := js.Global().Get("document").Call("createElement", "div")
		gridEl.Set("id", fmt.Sprintf("grid-layer-%d", i))
		gridEl.Set("className", "grid")

		gridEl.Set("innerHTML", grid)

		js.Global().Get("document").Call("getElementById", "editor").Call("appendChild", gridEl)
	}

	js.Global().Get("document").Get("documentElement").Get("style").Call("setProperty", "--col-number", len(block.Tiles))
}

type Button struct {
	buttonType string
	directions []string
}

func getButtons(this js.Value, args []js.Value) interface{} {
	resetButton := `<button onclick="createBlock()">Reset Block</button>`
	saveButton := `<button onclick="saveBlock()">Save Block</button>`
	saveAndResetSection := `<section class="flex-column">` + resetButton + saveButton + "</section>"

	collisionButton := fmt.Sprintf(`<button onclick="EDITOR.tileType='%s'">Add collision</button>`, Collision)
	removeButton := fmt.Sprintf(`<button onclick="EDITOR.tileType='%s'">Remove tile</button>`, Empty)

	buttonGroup := `<section class="flex-column">`

	for tileType, tiles := range getTileset() {
		buttonGroup += fmt.Sprintf(`<button onclick="showButtons('%s')">%s</button>`, tileType, tileType)
		buttonGroup += fmt.Sprintf(`<div id="bgroup-%s" class="bgroup hidden">`, tileType)

		for _, tile := range tiles {
			buttonGroup += fmt.Sprintf(`<button class="tile" style="background-position: %dpx %dpx" onclick="EDITOR.tileType='%s'; EDITOR.commandParams = {tileId: %d};"></button>`, tile.x*ButtonTileMultiplier, tile.y*ButtonTileMultiplier, tileType, tile.id)
		}

		buttonGroup += "</div>"
	}

	buttonGroup += "</section>"

	connections := `<section>Tile Connections<br>`
	connections += fmt.Sprintf(`<input type="checkbox" id="up" name="connections" value="up" %s><label for="up">top</label>`, conditionalAttribute(block.Connections[0], "checked"))
	connections += fmt.Sprintf(`<input type="checkbox" id="right" name="connections" value="right" %s><label for="right">right</label>`, conditionalAttribute(block.Connections[1], "checked"))
	connections += fmt.Sprintf(`<input type="checkbox" id="down" name="connections" value="down" %s><label for="down">bottom</label>`, conditionalAttribute(block.Connections[3], "checked"))
	connections += fmt.Sprintf(`<input type="checkbox" id="left" name="connections" value="left" %s><label for="left">left</label>`, conditionalAttribute(block.Connections[3], "checked"))
	connections += `</section>`

	otherButtons := `<section id="action-buttons"><button onclick="_EDITOR_changeZoom(1)"><i class="zoom-in"></i></button>`
	otherButtons += fmt.Sprintf(`<button onclick="_EDITOR_changeZoom(-1)" %s><i class="zoom-out"></i></button>`, conditionalAttribute(currentScale < 2, "disabled"))
	otherButtons += `<button onclick="toggleGrid()"><i class="display-grid"></i></button></section>`

	js.Global().Get("document").Call("getElementById", "command-panel").Set("innerHTML", saveAndResetSection+buttonGroup+collisionButton+removeButton+connections+otherButtons)

	return nil
}

func createBlock(this js.Value, args []js.Value) interface{} {
	gridSize := DefaultGridSize

	if len(args) > 0 {
		gridSize = args[0].Int()
	}

	// Create a new block instance
	block = Block{
		id:          "",
		Tiles:       make([][]Tile, gridSize),
		Connections: [4]bool{},
	}

	for i := range block.Tiles {
		block.Tiles[i] = make([]Tile, gridSize)
		for j := range block.Tiles[i] {
			layer := []Layer{{MaterialType: "floor", TileId: 150}}
			defaultTile := Tile{Layers: layer, IsBlocking: false}
			block.Tiles[i][j] = defaultTile
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
		fmt.Println("This tile type does not exist: ", tileType)
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
			fmt.Printf("Error marshaling block data: %v", err)
			return
		}

		httpMethod := fetch.MethodPost

		if block.id != "" {
			httpMethod = fetch.MethodPut
		}

		resp, _ := fetch.Fetch("http://localhost:8081/api/blocks/"+block.id, &fetch.Opts{
			Body:   bytes.NewBuffer(data),
			Method: httpMethod,
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

		fmt.Println("load block")

		block.id = blockId

		renderGrid()
		getButtons(this, args)
		renderLayersSection(this, args)
	}()
	return nil
}

// args[0] -> zoom 1 -> in | -1 -> out
func changeZoom(this js.Value, args []js.Value) interface{} {
	zoom := int16(-1)

	if !args[0].IsUndefined() && args[0].Int() == 1 {
		zoom = 1
	}

	currentScale += zoom

	js.Global().Get("document").Get("documentElement").Get("style").Call("setProperty", "--scale-value", currentScale)

	renderGrid()
	getButtons(this, args)
	renderLayersSection(this, args)

	return nil
}

func addLayer(this js.Value, args []js.Value) interface{} {
	for y := range block.Tiles {
		for x := range block.Tiles[y] {
			block.Tiles[y][x].Layers = append(block.Tiles[y][x].Layers, Layer{})
		}
	}

	renderLayersSection(this, args)

	return nil
}

func renderLayersSection(this js.Value, args []js.Value) interface{} {
	layers := len(block.Tiles[0][0].Layers)

	section := `
    <div class="layers">    
      <div id="layer-0">
        <span onclick="selectLayer()">Collisions</span> 
      </div>
  `

	for i := 1; i <= layers; i++ {
		section += fmt.Sprintf(`
      <div id="layer-%d"> 
        <span onclick="selectLayer(%d)">Layer %d</span> 
        <button><i>H</i></button>
        <button><i>D</i></button>
      </div>
      `, i, i-1, i)
	}

	// Collisions don't have a layer
	section += `
    </div>
      <div>
        <button onclick="_EDITOR_addLayer()">Add Layer</button>
      </div>`

	parent := js.Global().Get("document").Call("getElementById", "command-panel")

	currentLayerSection := js.Global().Get("document").Call("getElementById", "layers-section")

	newLayerSection := js.Global().Get("document").Call("createElement", "section")
	newLayerSection.Set("id", "layers-section")
	newLayerSection.Set("innerHTML", section)

	if currentLayerSection.Truthy() {
		parent.Call("replaceChild", newLayerSection, currentLayerSection)
	} else {
		parent.Call("appendChild", newLayerSection)
	}

	return nil
}
