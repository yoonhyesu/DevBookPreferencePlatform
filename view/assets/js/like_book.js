$(document).ready(function () {
    // 현재 사용자 정보 가져오기
    const currentUserElement = document.getElementById('current-user');

    // 사용자 요소가 없는 경우 처리
    if (!currentUserElement) {
        console.error('사용자 정보 요소를 찾을 수 없습니다');
        return;
    }

    const userId = currentUserElement.dataset.userId;

    // 로그인 상태 확인
    if (!userId) {
        console.warn('로그인되지 않은 상태입니다');
        window.location.href = '/login';
        return;
    }

    console.log('좋아요 목록 페이지 로드됨');

    // 좋아요 취소 버튼 클릭 이벤트
    $('.remove-like').on('click', function (e) {
        e.preventDefault();
        const bookId = $(this).data('book-id');
        const card = $(this).closest('.col');

        console.log('좋아요 취소 요청:', bookId);

        $.ajax({
            url: '/api/books/like',
            type: 'POST',
            data: {
                bookId: bookId,
                liked: false
            },
            success: function (response) {
                console.log('좋아요 취소 성공:', response);

                // 성공 시 카드 제거 (애니메이션 효과 추가)
                card.fadeOut(300, function () {
                    $(this).remove();

                    // 남은 책이 없으면 빈 상태 메시지 표시
                    if ($('.book-card').length === 0) {
                        $('.row-cols-2').html(`
                            <div class="text-center py-5 w-100">
                                <i class="bi bi-heart" style="font-size: 3rem;"></i>
                                <h5 class="mt-3">아직 좋아요한 책이 없습니다</h5>
                                <p class="text-muted">마음에 드는 책을 찾아 좋아요를 눌러보세요!</p>
                                <a href="/" class="btn btn-primary mt-2">책 둘러보기</a>
                            </div>
                        `);
                    }
                });
            },
            error: function (error) {
                console.error('좋아요 취소 중 오류 발생:', error);
                alert('좋아요 취소 중 오류가 발생했습니다. 다시 시도해주세요.');
            }
        });
    });

    // 이미지 로드 오류 처리
    $('.card-img-top').on('error', function () {
        // 이미지 로드 실패 시 기본 이미지로 대체
        $(this).attr('src', '/assets/images/book-placeholder.png');
        console.log('이미지 로드 실패, 기본 이미지로 대체:', $(this).attr('alt'));
    });

    // 책 카드 호버 효과
    $('.book-card').hover(
        function () {
            $(this).addClass('shadow-sm');
            $(this).css('transform', 'translateY(-5px)');
            $(this).css('transition', 'transform 0.3s ease');
        },
        function () {
            $(this).removeClass('shadow-sm');
            $(this).css('transform', 'translateY(0)');
        }
    );

    // 페이지 로드 시 책 표지 이미지 로드 확인
    console.log('좋아요한 책 수:', $('.book-card').length);
    $('.card-img-top').each(function () {
        console.log('책 표지 로드:', $(this).attr('src'));
    });
});
