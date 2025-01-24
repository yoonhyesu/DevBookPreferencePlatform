let ws;
let bookId, currentUserId, currentUserName;

function formatDate(date) {
  return new Date(date).toLocaleString('ko-KR', {
    hour: '2-digit',
    minute: '2-digit',
  });
}

function appendMessage(userName, side, text, time) {
  const chatContainer = document.querySelector('.chat-body');
  if (!chatContainer) {
    console.error('채팅 컨테이너를 찾을 수 없습니다');
    return;
  }

  const msgHTML = `
        <div class="msg ${side}-msg">
            <div class="msg-bubble ${side}-bubble">
                <div class="msg-info">
                    <div class="msg-info-name">${userId}(${userName})</div>
                    <div class="msg-info-time">${time}</div>
                </div>
                <div class="msg-text">${text}</div>
            </div>
        </div>
    `;

  chatContainer.insertAdjacentHTML("beforeend", msgHTML);
  chatContainer.scrollTop = chatContainer.scrollHeight;
}

function connectWebSocket(bookId) {
  if (!bookId) {
    console.error('책 ID가 없습니다');
    return null;
  }

  // 현재 호스트 기반으로 웹소켓 주소 생성
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  const ws = new WebSocket(`${protocol}//${host}/ws/chat?BOOK_ID=${bookId}`);

  // 재연결 시도 횟수 제한
  let retryCount = 0;
  const maxRetries = 3;

  ws.onopen = function () {
    console.log("웹소켓 연결 성공!");
    retryCount = 0; // 연결 성공시 카운트 리셋
  };

  ws.onmessage = function (event) {
    try {
      const message = JSON.parse(event.data);
      displayMessage(message);
    } catch (error) {
      console.error("메시지 처리 실패:", error);
    }
  };

  ws.onerror = function (error) {
    console.error("웹소켓 에러:", error);
    if (retryCount < maxRetries) {
      retryCount++;
      console.log(`재연결 시도 ${retryCount}/${maxRetries}...`);
      setTimeout(() => connectWebSocket(bookId), 3000);
    } else {
      console.log("최대 재시도 횟수 초과");
    }
  };

  ws.onclose = function () {
    console.log("웹소켓 연결 종료");
  };

  return ws;
}

function displayMessage(message) {
  const chatContainer = document.querySelector('.chat-body');
  if (!chatContainer) {
    console.error('채팅 컨테이너를 찾을 수 없습니다');
    return;
  }

  const messageDiv = document.createElement('div');
  messageDiv.className = `msg ${message.USER_ID === currentUserId ? 'right' : 'left'}-msg`;

  messageDiv.innerHTML = `
    <div class="msg-bubble ${message.USER_ID === currentUserId ? 'right' : 'left'}-bubble">
        <div class="msg-info">
            <div class="msg-info-name">${message.USER_NAME}</div>
            <div class="msg-info-time">${new Date(message.CREATED_AT).toLocaleTimeString()}</div>
        </div>
        <div class="msg-text">${message.MESSAGE}</div>
    </div>
  `;

  chatContainer.appendChild(messageDiv);
  chatContainer.scrollTop = chatContainer.scrollHeight;
}

// 페이지 로드 시 실행
document.addEventListener('DOMContentLoaded', () => {
  // DOM 요소 확인 후 값 설정
  const bookTitleElement = document.getElementById('book-title');
  const currentUserElement = document.getElementById('current-user');

  if (!bookTitleElement || !currentUserElement) {
    console.error('필요한 DOM 요소를 찾을 수 없습니다:', {
      bookTitle: bookTitleElement,
      currentUser: currentUserElement
    });
    return;
  }

  // 데이터 확인 로깅
  console.log('Book ID:', bookTitleElement.dataset.bookId);
  console.log('User Info:', {
    id: currentUserElement.dataset.userId,
    name: currentUserElement.dataset.userName
  });

  // 데이터 설정
  bookId = bookTitleElement.dataset.bookId;
  currentUserId = currentUserElement.dataset.userId;
  currentUserName = currentUserElement.dataset.userName;

  if (!bookId || !currentUserId || !currentUserName) {
    console.error('필수 데이터가 없습니다:', {
      bookId,
      currentUserId,
      currentUserName
    });
    return;
  }

  // 채팅 초기화
  initializeChat();
  console.log("작동하냐?")
});

function initializeChat() {
  const msgForm = document.getElementById('msgForm');
  const msgInput = document.getElementById('msgInput');
  const chatContainer = document.querySelector('.chat-body');

  if (!msgForm || !msgInput || !chatContainer) {
    console.error('채팅 UI 요소를 찾을 수 없습니다:', {
      form: msgForm,
      input: msgInput,
      container: chatContainer
    });
    return;
  }

  // 웹소켓 연결
  ws = connectWebSocket(bookId);

  // 이전 메시지 로드
  loadPreviousMessages();

  // 메시지 전송 이벤트 리스너
  msgForm.addEventListener('submit', (event) => {
    event.preventDefault();
    const msgText = msgInput.value.trim();

    // 웹소켓 상태 확인
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      console.error('웹소켓이 연결되어 있지 않습니다');
      alert('채팅 서버와 연결이 끊어졌습니다. 페이지를 새로고침해주세요.');
      return;
    }

    if (!msgText) return;

    const message = {
      BOOK_ID: parseInt(bookId),
      USER_ID: currentUserId,
      USER_NAME: currentUserName,
      MESSAGE: msgText
    };

    try {
      ws.send(JSON.stringify(message));
      msgInput.value = "";
    } catch (error) {
      console.error('메시지 전송 실패:', error);
      alert('메시지 전송에 실패했습니다.');
    }
  });
}

// 이전 메시지 로드
function loadPreviousMessages() {
  if (!bookId) {
    console.error('책 ID가 없어 이전 메시지를 로드할 수 없습니다');
    return;
  }
  console.log('bookid', bookId)
  fetch(`/api/chat/messages/${bookId}`)
    .then(response => {
      if (!response.ok) {
        throw new Error('이전 메시지 로드 실패');
      }
      return response.json();
    })
    .then(messages => {
      if (Array.isArray(messages)) {
        messages.reverse().forEach(msg => {
          const side = msg.USER_ID === currentUserId ? 'right' : 'left';
          appendMessage(msg.USER_NAME, side, msg.MESSAGE, formatDate(msg.CREATED_AT));
        });
      }
    })
    .catch(error => console.error('이전 메시지 로드 실패:', error));
}

// 이전 메시지 로드

