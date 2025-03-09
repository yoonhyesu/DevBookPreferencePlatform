$(function () {
    const bookId = $('#book-title').data('book-id');
    $(".chat-menu-icons .toogle-bar").click(function () {
        $(".chat-menu").toggleClass("show");
    });
    GetContentList(bookId);
});

// 좋아요 버튼 기능 구현
$(document).ready(function () {
    const likeButton = $('#likeButton');
    const currentUserElement = document.getElementById('current-user');

    // 사용자 요소가 없는 경우 처리
    if (!currentUserElement) {
        console.error('사용자 정보 요소를 찾을 수 없습니다');
        return;
    }

    const userId = currentUserElement.dataset.userId;

    console.log('User Info:', {
        id: currentUserElement.dataset.userId,
        name: currentUserElement.dataset.userName
    });

    // 로그인 상태 확인
    if (!userId) {
        console.warn('로그인되지 않은 상태입니다');
        // 좋아요 버튼 비활성화 또는 로그인 유도 UI 표시 가능
    }

    // 좋아요 버튼 클릭 이벤트
    likeButton.on('click', function () {
        // 로그인 확인
        if (!userId) {
            alert('로그인이 필요한 서비스입니다.');
            window.location.href = '/login';
            return;
        }

        const icon = $(this).find('i');
        const bookId = $(this).data('book-id');

        // 아이콘 클래스 토글 (빈 하트 <-> 채워진 하트)
        let isLiked = false;
        if (icon.hasClass('bi-heart')) {
            icon.removeClass('bi-heart').addClass('bi-heart-fill text-danger');
            isLiked = true;
        } else {
            icon.removeClass('bi-heart-fill text-danger').addClass('bi-heart');
            isLiked = false;
        }

        // 서버로 좋아요 상태 저장 요청
        $.ajax({
            url: '/api/books/like',
            type: 'POST',
            data: {
                bookId: bookId,
                liked: isLiked
            },
            success: function (response) {
                console.log('좋아요 상태가 저장되었습니다.');
            },
            error: function (error) {
                console.error('좋아요 저장 중 오류 발생:', error);
                // 오류 발생 시 아이콘 상태 되돌리기
                if (isLiked) {
                    icon.removeClass('bi-heart-fill text-danger').addClass('bi-heart');
                } else {
                    icon.removeClass('bi-heart').addClass('bi-heart-fill text-danger');
                }
            }
        });
    });

    // 페이지 로드 시 사용자의 좋아요 상태 확인
    if (userId) {
        const bookId = likeButton.data('book-id');
        $.ajax({
            url: '/api/books/like/status',
            type: 'GET',
            data: {
                bookId: bookId
            },
            success: function (response) {
                if (response.liked) {
                    likeButton.find('i').removeClass('bi-heart').addClass('bi-heart-fill text-danger');
                }
            }
        });
    }
});

// 목차 바인딩 프로세스
function GetContentList(bookId) {
    console.log("목차 로딩 시작");
    $.ajax({
        url: '/book/contents/' + bookId,
        method: 'GET',
        success: function (response) {
            console.log("책 목차 로드 성공:", response);
            var result = response.contents;
            var contentList = result
                .replace(/\\r/g, '')
                .replace(/\\n/g, '<br><br>')
                .replace(/\\t/g, '')
                .replace(/tr$/gm, '')
                .replace(/<br><br>/g, '<br>');

            // 목차를 div로 감싸서 스타일 적용
            $('#contents-list').html('<div class="contents-text">' + contentList + '</div>');

            // 토글 버튼 추가
            if ($('.toggle-button-wrapper').length === 0) {
                $('#contents-list').after(`
                    <div class="toggle-button-wrapper text-center mt-2">
                        <button class="btn btn-light btn-sm toggle-contents">
                            <span class="expand-text">펼쳐보기</span>
                            <span class="collapse-text" style="display:none;">접어보기</span>
                            <i class="bi bi-chevron-down expand-icon"></i>
                            <i class="bi bi-chevron-up collapse-icon" style="display:none;"></i>
                        </button>
                    </div>
                `);
            }

            // 목차 높이 제한 및 토글 기능
            const contentsText = $('.contents-text');
            const originalHeight = contentsText.height();
            const maxHeight = 200; // 최대 높이 (px)

            if (originalHeight > maxHeight) {
                // 목차가 길면 높이 제한 및 토글 버튼 표시
                contentsText.css({
                    'max-height': maxHeight + 'px',
                    'overflow': 'hidden'
                });

                $('.toggle-contents').show().on('click', function () {
                    if (contentsText.css('max-height') !== 'none') {
                        // 펼치기
                        contentsText.css('max-height', 'none');
                        $('.expand-text, .expand-icon').hide();
                        $('.collapse-text, .collapse-icon').show();
                    } else {
                        // 접기
                        contentsText.css('max-height', maxHeight + 'px');
                        $('.expand-text, .expand-icon').show();
                        $('.collapse-text, .collapse-icon').hide();
                    }
                });
            } else {
                // 목차가 짧으면 토글 버튼 숨김
                $('.toggle-button-wrapper').hide();
            }
        },
        error: function (error) {
            console.error("목차 로드 실패:", error);
            $('#contents-list').html('<p class="text-muted">목차 정보를 불러올 수 없습니다.</p>');
        }
    });
}