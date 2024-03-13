const EDITOR = {
  tileType: "",
  commandParams: {},
  layer: 0,
};

function createBlock() {
  _EDITOR_createBlock();
  _EDITOR_renderSidenav();
  _EDITOR_renderLayersSection();
}

function toggleGrid() {
  document.getElementsByClassName("grid")[0].classList.toggle("no-borders");
}

function showButtons(index) {
  document.getElementById("bgroup-" + index).classList.toggle("hidden");
}

function selectLayer(index) {
  EDITOR.layer = index;
  const layers = document.querySelectorAll("#layers-section > .layers > div");
  layers.forEach((b) => b.classList.remove("active"));
  layers[(index ?? -1) + 1].classList.add("active");

  const collisions = document.querySelectorAll(".collision");
  collisions.forEach((c) => c.classList.remove("show"));

  if (index === undefined) {
    collisions.forEach((div) => div.classList.add("show"));
  }
}

function saveBlock() {
  const directions = [false, false, false, false];
  const directionElements = document.querySelectorAll(
    "input[type='checkbox'][name='connections']",
  );

  for (let i = 0; i < 4; i++) {
    directions[i] = directionElements[i].checked;
  }

  // TODO: might be needed in the future
  EDITOR.commandParams = { directions: directions };

  _EDITOR_saveBlock(directions);
}

function setTile(coordinates) {
  if (!EDITOR.tileType) {
    return;
  }

  _EDITOR_setTile(
    coordinates,
    EDITOR.tileType,
    EDITOR.commandParams.tileId,
    EDITOR.layer,
  );

  // If collision, select the collision layer
  if (EDITOR.tileType === "collision") {
    selectLayer();
  } else {
    selectLayer(EDITOR.layer);
  }
}

function toggleLayerVisibility(index, el) {
  document.querySelector(`#grid-layer-${index}`)?.classList.toggle("hidden");
  el.classList.toggle("active");
}
