const messages = document.getElementById('messages');
const usernameInput = document.getElementById('username');
const messageInput = document.getElementById('message');
const sendButton = document.getElementById('send');

let ws;

// 初始化WebSocket连接
function connect() {
    ws = new WebSocket('ws://' + window.location.host + '/ws');
    ws.onmessage = function(event) {
        const msg = JSON.parse(event.data);
        const messageElement = document.createElement('div');
        messageElement.className = 'message';
        messageElement.innerHTML = `<span class="username">${msg.username}</span>: ${msg.content}`;
        messages.appendChild(messageElement);
        messages.scrollTop = messages.scrollHeight;
    };
    ws.onclose = function(event) {
        console.log('WebSocket closed', event);
    };
}

// 发送消息
function sendMessage() {
    const username = usernameInput.value.trim();
    const message = messageInput.value.trim();
    if (username && message) {
        ws.send(JSON.stringify({ username, content: message }));
        messageInput.value = '';
    }
}

// 绑定事件
window.onload = function () {
    connect();
    sendButton.addEventListener('click', sendMessage);
};
