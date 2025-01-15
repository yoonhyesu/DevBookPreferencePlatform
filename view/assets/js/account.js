$(document).ready(function () {
    let flagPW = false;
});

/* 현재 비밀번호 확인 */
function checkPassword() {
    var current_pw = $('#current_pw').val();
    console.log("입력된 비밀번호 길이:", current_pw.length);

    if (current_pw.length >= 7) {
        $.ajax({
            url: "/account/check_pw",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                CURRENT_PASSWORD: current_pw
            }),
            success: function (response) {
                console.log("서버 응답:", response);
                if (response.success) {
                    $('#current_pw_check').text('비밀번호 일치').css('color', 'green');
                    flagPW = true;
                } else {
                    $('#current_pw_check').text('비밀번호 불일치').css('color', 'red');
                    flagPW = false;
                }
            },
            error: function (xhr, status, error) {
                console.error("비밀번호 확인 중 오류 발생:", error);
                console.error("서버 응답:", xhr.responseText);
                $('#current_pw_check').text('비밀번호 불일치').css('color', 'red');
                flagPW = false;
            }
        });
    }
}

/* 신규 비밀번호 유효성 검사 */
function validateNewPassword() {
    const regexPw = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,16}$/;
    var current_pw = $('#current_pw').val();
    checkPW = $('#new_pw').val();

    if (current_pw === checkPW) {
        $('#new_pw_message').text('현재 비밀번호와 동일합니다.').css('color', 'red');
        flagPW = false;
    } else if (!regexPw.test(checkPW)) {
        $('#new_pw_message').text('비밀번호는 8~16자의 영문 대/소문자, 숫자, 특수문자(!,@,#,$,%,^,&,*,-)를 사용해주세요.').css('color', 'red')
        flagPW = false;
    } else {
        $('#new_pw_message').text('사용 가능한 비밀번호 형식입니다.').css('color', 'green');
        flagPW = true;
    }
}


/* 신규 비밀번호 재확인 */
function checkNewPassword() {
    var new_pw = $('#new_pw').val();
    var new_pw_check = $('#new_pw_check').val();
    if (new_pw === new_pw_check) {
        $('#new_pw_check').text('비밀번호 일치').css('color', 'green');
        $('#new_pw_check_message').text('비밀번호 일치').css('color', 'green');
        flagPW = true;
    } else {
        $('#new_pw_check').text('비밀번호 불일치').css('color', 'red');
        $('#new_pw_check_message').text('비밀번호 불일치').css('color', 'red');
        flagPW = false;
    }
}

/* 비밀번호 변경 */
$('#pw_change_btn').on('click', function (e) {
    e.preventDefault();
    if (flagPW === true) {
        $.ajax({
            url: '/account/profile-setting/changpw',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                USER_ID: $('#user_id').val(),
                OLD_PASSWORD: $('#current_pw').val(),
                NEW_PASSWORD: $('#new_pw').val()
            })
        }).done(function (response) {
            alert('비밀번호가 성공적으로 변경되었습니다!!!');
            location.reload();
        }).fail(function (error) {
            console.error('비밀번호 변경 실패:', error);
            alert(error.responseJSON?.error || '비밀번호 변경에 실패했습니다!!!');
        });
    }
});

/* 회원 삭제 */
$('#deleteUser_btn').on('click', function (event) {
    event.preventDefault();
    if ($('#delete_check').is(':checked')) {
        $.ajax({
            url: '/account/profile-setting/leave',
            method: 'POST',
            dataType: 'json',
            contentType: 'application/json',
            data: JSON.stringify({})
        })
            .done(function (response) {
                alert("회원탈퇴가 완료되었습니다!");
                window.location.href = "/";
            })
            .fail(function (jqXHR, textStatus, errorThrown) {
                console.error("AJAX 오류:", textStatus, errorThrown);
                console.log("서버 응답:", jqXHR.responseText);
                alert(jqXHR.responseJSON?.error || '회원 삭제 중 오류가 발생했습니다. 다시 시도해주세요.');
            });
    } else {
        alert('회원 탈퇴 동의에 체크해주세요.');
    }
});


