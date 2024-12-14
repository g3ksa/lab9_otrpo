let ws
let nameInput = document.getElementById('nameInput')
let connectBtn = document.getElementById('connectBtn')
let usersList = document.getElementById('users')
let messagesDiv = document.getElementById('messages')
let msgInput = document.getElementById('msgInput')
let sendBtn = document.getElementById('sendBtn')

connectBtn.onclick = () => {
	const name = nameInput.value.trim()
	if (!name) {
		alert('Введите имя!')
		return
	}
	ws = new WebSocket(
		`ws://${location.host}/ws?name=${encodeURIComponent(name)}`
	)

	ws.onopen = () => {
		console.log('WS Connected')
	}

	ws.onmessage = (evt) => {
		const msg = JSON.parse(evt.data)
		if (msg.type === 'message') {
			appendMessage(`${msg.author}: ${msg.text}`)
		} else if (msg.type === 'users') {
			usersList.innerHTML = ''
			msg.users.forEach((u) => {
				let li = document.createElement('li')
				li.textContent = u
				usersList.appendChild(li)
			})
		}
	}

	ws.onclose = () => {
		appendMessage('Подключение закрыто')
	}
}

sendBtn.onclick = () => {
	if (ws && ws.readyState === WebSocket.OPEN) {
		ws.send(msgInput.value)
		msgInput.value = ''
		msgInput.focus()
	} else {
		alert('Не подключены к чату.')
	}
}

function appendMessage(text) {
	let div = document.createElement('div')
	div.textContent = text
	messagesDiv.appendChild(div)
	messagesDiv.scrollTop = messagesDiv.scrollHeight
}
