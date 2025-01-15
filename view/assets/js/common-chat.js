let ws;
const bookId = document.getElementById('book-title').dataset.bookId;
const currentUserId = document.getElementById('current-user').dataset.userId;
const currentUserName = document.getElementById('current-user').dataset.userName;
const msgForm = document.getElementById('msgForm');
const msgInput = document.getElementById('msgInput');
const msgerChat = document.getElementById('msgChat');

function formatDate(date) {
  return new Date(date).toLocaleString('ko-KR', {
    hour: '2-digit',
    minute: '2-digit',
  });
}

function appendMessage(userName, side, text, time) {
  const msgHTML = `
        <div class="msg ${side}-msg">
            <div class="msg-img">
                <img src="/assets/images/profile.PNG" alt="Profile Image">
            </div>
            <div class="msg-bubble ${side}-bubble">
                <div class="msg-info">
                    <div class="msg-info-name">${userName}</div>
                    <div class="msg-info-time">${time}</div>
                </div>
                <div class="msg-text">${text}</div>
            </div>
        </div>
    `;

  msgerChat.insertAdjacentHTML("beforeend", msgHTML);
  msgerChat.scrollTop = msgerChat.scrollHeight;
}

function connectWebSocket() {
  ws = new WebSocket(`ws://${window.location.host}/ws/chat/${bookId}`);

  ws.onmessage = function (event) {
    const msg = JSON.parse(event.data);
    const side = msg.REG_USER_ID === currentUserId ? 'right' : 'left';
    appendMessage(msg.REG_USER_NAME, side, msg.MESSAGE, formatDate(msg.REG_DATE));
  };

  ws.onerror = function (error) {
    console.error('WebSocket 에러:', error);
  };

  ws.onclose = function () {
    console.log('WebSocket 연결 종료');
    setTimeout(connectWebSocket, 3000);
  };
}

msgForm.addEventListener('submit', (event) => {
  event.preventDefault();
  const msgText = msgInput.value.trim();
  if (!msgText) return;

  const message = {
    BOOK_ID: parseInt(bookId),
    REG_USER_ID: currentUserId,
    REG_USER_NAME: currentUserName,
    MESSAGE: msgText
  };

  ws.send(JSON.stringify(message));
  msgInput.value = "";
});

// 이전 메시지 로드
function loadPreviousMessages() {
  fetch(`/api/chat/messages/${bookId}`)
    .then(response => response.json())
    .then(messages => {
      messages.reverse().forEach(msg => {
        const side = msg.REG_USER_ID === currentUserId ? 'right' : 'left';
        appendMessage(msg.REG_USER_NAME, side, msg.MESSAGE, formatDate(msg.REG_DATE));
      });
    })
    .catch(error => console.error('이전 메시지 로드 실패:', error));
}

// 페이지 로드시 실행
document.addEventListener('DOMContentLoaded', () => {
  connectWebSocket();
  loadPreviousMessages();
});
