const EDITOR = {
  commandType: "",
  commandParams: {},
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

function selectTile(coordinates) {
  if (!EDITOR.commandType) {
    return;
  }

  _EDITOR_selectTile(
    coordinates,
    EDITOR.commandType,
    EDITOR.commandParams.isBlocking || false,
    EDITOR.commandParams.direction || "top",
  );
}
