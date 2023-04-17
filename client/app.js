// creating the websocket connection
var socket = new WebSocket("wss://waffle-game-theory.onrender.com/ws");

// creating a unique id for player to know and send
var id = Date.now();

// when an update is received via ws connection, we update the model
socket.onmessage = function (evt) {
  var data = JSON.parse(evt.data);
  if (data.reset) {
    reset();
  }
  interpretInfo(data);
  game.gameState = data;
};

// break down information
function interpretInfo(data) {
  if (data.results.playerOne) {
    game.opponentReady = true;
    resolveWinner(data);
  } else if (data.playerOne.status && !game.playerReady) {
    game.opponentReady = true;
  }
}

// find and display winner
function resolveWinner(data) {
  var same;
  var playerOneResult = data.results.playerOne;
  var playerTwoResult = data.results.playerTwo;
  if (playerOneResult == playerTwoResult) {
    same = true;
    switch (playerOneResult) {
      case "split":
        game.winnerMessage = "Waffles for everybody!!";
        break;
      case "take":
        game.winnerMessage = "Oh no!! The waffle is trashed!";
    }
  }
  if (id == data.playerOne.id) {
    if (!same && data.results.playerOne == "split") {
      game.winnerMessage = "Your friend ate the waffle";
    } else if (!same && data.results.playerOne == "take") {
      game.winnerMessage = "You ate the whole waffle!";
    }
    game.playerChoice = data.results.playerOne;
    game.opponentChoice = data.results.playerTwo;
  } else {
    if (!same && data.results.playerTwo == "split") {
      game.winnerMessage = "Your friend took ate the waffle";
    } else if (!same && data.results.playerTwo == "take") {
      game.winnerMessage = "You eat the whole waffle!";
    }
    game.playerChoice = data.results.playerTwo;
    game.opponentChoice = data.results.playerOne;
  }
}

// reset the game in command
function reset() {
  game.playerReady = false;
  game.opponentReady = false;
  game.winnerMessage = "";
  game.playerChoice = "";
  game.opponentChoice = "";
}

// vuejs debug mode
Vue.config.debug = true; //TODO: Remove in production

// transistions
Vue.transition("board", {
  enterClass: "bounceInDown",
  leaveClass: "bounceOutDown",
});

// creating the vue instance to send information
var game = new Vue({
  el: "#game",
  data: {
    gameState: {
      started: false,
      messages: [],
    },
    text: "",
    choice: "",
    playerReady: "",
    opponentReady: "",
    winnerMessage: "",
    playerChoice: "",
    opponentChoice: "",
  },
  methods: {
    submit: function () {
      this.playerReady = true;
      var message = JSON.stringify({
        type: "submission",
        data: this.choice,
        id: id,
      });
      socket.send(message);
    },
    message: function () {
      if (this.text != "") {
        var message = JSON.stringify({
          type: "message",
          data: this.text,
          id: id,
        });
        socket.send(message);
        this.text = "";
      }
    },
    reset: function () {
      var message = JSON.stringify({ type: "command", data: "reset", id: id });
      socket.send(message);
    },
  },
});
