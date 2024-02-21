# RoK Dungeon

A dungeon crawler game based on Raking of Kings (王様ランキング) made with Golang, WebAssembly and JavaScript.

The game features an editor that empowers you to shape the very tower Bojji ascends. Populate ascending floors with creatures, set traps, and script events to create challenges that reflect the essence of Bojji's ascent. Kage, your ever-helpful companion, will assist in unraveling the secrets hidden within the ascending tower, providing guidance and aid as you climb to new heights.

(Thanks ChatGPT)

## Build and Run

If you have [goexec](https://github.com/shurcooL/goexec) installed, simply run

```
bash ./build.sh
```

Or you can build it and use your favorite web server with

```
GOOS=js GOARCH=wasm go build -o dist/editor.wasm game-engine/editor.go game-engine/utils.go
```

Editor -> localhost:8080/web/editor.html

## Features

### Editor

- [ ] Tileset and Map Editor
  - [x] Add collisions
  - [ ] Add layers
  - [ ] Add different types of walls, grounds, etc
  - [ ] Save current block
  - [ ] Load current block
  - [ ] Choose a size for the block
  - [ ] Floor Misc
  - [ ] Wall Misc
- [ ] Procedural Generation Controls
- [ ] Enemy Placement and Behavior
- [ ] Item and Loot Editor
- [ ] Event Trigger System
- [ ] Group floors and tiles by complexity and difficulty
- [ ] Dialogue and Quest Editor
- [ ] Undo/Redo Functionality

### Game

- [ ] Procedurally Generated floors
- [ ] Random events
- [ ] Boss battles
- [ ] Character Progression
- [ ] Items, weapons and armorms
- [ ] Tactical combat
- [ ] Sound and Music
- [ ] Level progression
- [ ] Multi language support (ENG and JPN)

## TODO

- [ ] Add a proper console log for wasm
- [ ] Add JSdoc for the client
- [ ] Make a proper build
- [ ] Error handling
