let data = {
  showIdentity: true,
  showMain: false,
  rolls: [],
  who: null,
}

let app = null
let ws = null

function handleMessage(message) {
  console.log("incoming message:", message)

  if ("type" in message) {
    if (message.type == "roll") {
      data.rolls.push({
        who: message.message.who,
        values: message.message.rolls,
        result: message.message.result,
      })
    }
  }
}

function setupRollerForm(selector) {
  const form = document.querySelector(selector)
  form.addEventListener("submit", event => {
    event.preventDefault()

    if (ws == null) {
      alert("Connection is closed :(")
      return
    }
    ws.send(JSON.stringify({
      type: "roll",
      message: {
        who: data.who,
      },
    }))
  })
}

function setupLoginForm(selector) {
  const form = document.querySelector(selector)
  if (form === null) {
    return
  }

  const usernameField = form.querySelector("#username")
  form.addEventListener("submit", event => {
    event.preventDefault()

    let username = usernameField.value
    if (username === "") {
      return
    }

    data.who = username
    data.showMain = true
    data.showIdentity = false
  })
}

Vue.component("die", {
  props: ['value'],
  template: '<li class="die"><p>{{ value }}</p></li>',
})

Vue.component("dice-roll", {
  props: ['roll'],
  template: `
    <div class="roll">
      <strong>{{ roll.who }}</strong> rolled ({{ roll.result }}):
      <ul class="dice">
        <die v-for="value in roll.values" v-bind:value="value"></die>
      </ul>
    </div>
  `,
  mounted: function() {
    // console.log("created dice roll", this.$el)
    window.scrollTo(0, this.$el.offsetTop)
  }
})

function establishConnection() {
  const wsProto = (() => {
	if (location.protocol == "https:") {
	  return "wss://"
	}
	return "ws://"
  })()
  const wsAddr = `${wsProto}${location.host}/websocket`
  ws = new WebSocket(wsAddr)

  ws.addEventListener("message", event => {
    handleMessage(JSON.parse(event.data))
  })

  ws.addEventListener("close", event => {
    setTimeout(() => {
      establishConnection()
    }, 250)
  })

}

document.addEventListener("DOMContentLoaded", event => {
  establishConnection()

  app = new Vue({
    el: "#app",
    data: data,
  })

  setupLoginForm("#loginForm")
  setupRollerForm("#rollForm")
})
