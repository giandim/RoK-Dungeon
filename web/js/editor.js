const EDITOR = {
  tileType: "",
  commandParams: {},
  layer: 0,
};

function createBlock() {
  _EDITOR_createBlock();
  _EDITOR_getButtons();
}

function toggleGrid() {
  document.getElementsByClassName("grid")[0].classList.toggle("no-gap");
}

function showButtons(index) {
  document.getElementById("bgroup-" + index).classList.toggle("hidden");
}

function selectLayer(index) {
  EDITOR.layer = index;
  const buttons = document.querySelectorAll(".layer-buttons > button");
  buttons.forEach((b) => b.classList.remove("active"));
  buttons[index ?? 2].classList.add("active");

  const collisions = document.querySelectorAll(".collision");
  collisions.forEach((c) => c.classList.remove("show"));

  if (index === undefined) {
    collisions.forEach((div) => div.classList.add("show"));
  }
}

function selectTile(coordinates) {
  if (!EDITOR.tileType) {
    return;
  }

  _EDITOR_setTile(
    coordinates,
    EDITOR.tileType,
    EDITOR.commandParams.direction || "top",
    EDITOR.layer,
  );

  // If collision, select the collision layer
  if (EDITOR.tileType === "collision") {
    selectLayer();
  } else {
    selectLayer(EDITOR.layer);
  }
}
