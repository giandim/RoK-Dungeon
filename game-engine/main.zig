const std = @import("std");

extern fn consoleLog(arg: i32) void;

export fn add(a: i32, b: i32) i32 {
    consoleLog(a + b);
    return a + b;
}

const Block = struct { id: []u8, tilesetId: u32, tiles: [][]Tile, connections: [4]bool };
const Tile = struct { isBlocking: bool, layers: []Layer };
const Layer = struct { materialType: []u8, tileId: u32 };

export fn generateFloor(blocks: [*]i32) void {
    consoleLog(blocks);
}
