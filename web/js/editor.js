const EDITOR = {
  commandType: "",
  commandParams: {},
};

function createBlock() {
  _EDITOR_createBlock();
  const buttons = _EDITOR_getButtons();
  console.log(buttons);
  renderButtons(buttons);
  // const grid = _ENGINE_selectTile("0, 0");
  //console.log(grid);
}

function renderButtons(buttons) {
  resetButton = `<button onclick="createBlock()">Reset Block</button>`;
  let buttonGroup = "";

  buttons.forEach((button, index) => {
    buttonGroup += `
      <button onclick="showButtons(${index})">${button.type}</button>
      <div id="bgroup-${index}" class="bgroup">
    `;

    if (!button.directions.length) {
      buttonGroup += `<button class="${button.type}" onclick="EDITOR.commandType='${button.type}'; EDITOR.commandParams = [];"></button>`;
    } else {
      button.directions.forEach((direction) => {
        buttonGroup += `<button class="${button.type} ${direction}" onclick="EDITOR.commandType='${button.type}'; EDITOR.commandParams = {direction: '${direction}'};"></button>`;
      });
    }

    buttonGroup += `</div>`;
  });

  document.getElementById("command-panel").innerHTML =
    resetButton + buttonGroup;
}

function selectTile(coordinates) {
  console.log(EDITOR);
  _EDITOR_selectTile(
    coordinates,
    EDITOR.commandType,
    EDITOR.commandParams.isBlocking || false,
    EDITOR.commandParams.direction || "top",
  );
}
