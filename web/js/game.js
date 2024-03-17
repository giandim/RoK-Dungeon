var memory = new WebAssembly.Memory({
  initial: 2 /* pages */,
  maximum: 2 /* pages */,
});

var importObject = {
  env: {
    consoleLog: (arg) => console.log(arg), // Useful for debugging on zig's side
    memory: memory,
  },
};

WebAssembly.instantiateStreaming(
  fetch("/zig-out/lib/game.wasm"),
  importObject,
).then((result) => {
  // const wasmMemoryArray = new Uint8Array(memory.buffer);
  console.log(result.instance.exports.add(2, 3));
  console.log(memory.buffer);
  generateFloor();
});

function generateFloor() {
  const blocks = [];

  fetch("/game-engine/assets/blocks/")
    .then((r) => r.text())
    .then((files) => {
      const parser = new DOMParser();
      const doc = parser.parseFromString(files, "text/html");
      const links = doc.querySelectorAll("a[href]");
      const fileList = Array.from(links)
        .map((link) => link.textContent.includes(".bin") && link.textContent)
        .filter((file) => file);

      const promises = fileList.map((file) =>
        fetch("/game-engine/assets/blocks/" + file).then((r) => r.text()),
      );

      Promise.all(promises).then((data) => {
        blocks.push(data);
        console.log(blocks);
      });
      console.log("File list:", fileList);
    });
}
