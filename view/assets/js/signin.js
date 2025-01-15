$(document).ready(function () {
    $('#loginBtn').on('click', function (e) {
        e.preventDefault();
        const userId = $('#InputID').val();
        const password = $('#InputPW').val();

        if (!userId || !password) {
            alert('아이디와 비밀번호를 모두 입력해주세요!!');
            return;
        }

        $.ajax({
            url: '/auth/signin',
            method: 'POST',
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
            xhrFields: {
                withCredentials: true
            },
            crossDomain: true,
            data: JSON.stringify({
                USER_ID: userId,
                PASSWORD: password,
            }),
            success: function (response) {
                console.log("서버 응답:", response);
                if (response.message === "로그인 성공") {
                    alert('로그인이 완료되었습니다.');
                    window.location.href = '/';
                }
            },
            error: function (xhr, status, error) {
                if (xhr.status === 401) {
                    alert('아이디 또는 비밀번호가 올바르지 않습니다!!');
                } else if (xhr.status === 500) {
                    alert('서버 오류가 발생했습니다. 잠시 후 다시 시도해주세요!!');
                }
                console.error("AJAX 오류:", status, error);
            }
        });
    });

    $('#logoutBtn').on('click', function (e) {
        $.ajax({
            url: '/auth/signout',
            method: 'POST',
            xhrFields: {
                withCredentials: true
            },
            crossDomain: true,
            success: function (response) {
                alert('로그아웃되었습니다');
                window.location.href = '/';
            },
            error: function (xhr, status, error) {
                alert('로그아웃 중 오류가 발생했습니다');
                console.error('로그아웃 오류:', error);
            }
        });
    });
});
