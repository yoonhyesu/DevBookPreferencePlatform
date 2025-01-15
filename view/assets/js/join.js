$(document).ready(function () {
    let checkID = null;
    let checkPW = null;
    let checkEmail = null;
    let checkPhone = null;
    let checkName = null;

    let flagID = false;
    let flagPW = false;
    let flagEmail = false;
    let flagPhone = false;
    let flagName = false;

    // 아이디 작성 시 
    $('#user_id').on("input", function () {
        validateInputId();
    });

    // 비밀번호 작성 시
    $('#password').on("input", function () {
        validateInputPw();
    })

    // 이름 작성 시
    $('#user_name').on("input", function () {
        validateInputName();
    })

    // 이메일 작성 시
    $('#email').on("input", function () {
        validateInputEmail();
    })

    // 전화번호 작성 시
    $('#phone_number').on("input", function () {
        validateInputPhone();
    })

    // 아이디 유효성 검사
    function validateInputId() {
        const regexId = /^[a-zA-Z0-9_-]{6,20}$/;
        checkID = $('#user_id').val();

        if (checkID === "") {
            $('#id-message').html('<i class="fa-solid fa-x"></i> 아이디를 입력해주세요.');
            $('#id-message').css('color', 'red');
            flagID = false;
            return;
        }

        if (!regexId.test(checkID)) {
            $('#id-message').html('<i class="fa-solid fa-x"></i> 아이디는 6~20자의 영문 대소문자, 숫자, 특수문자(-, _)만 사용 가능합니다.');
            $('#id-message').css('color', 'red');
            flagID = false;
        } else {
            $('#id-message').html('<i class="fa-solid fa-check"></i> 사용 가능한 아이디 형식입니다.');
            $('#id-message').css('color', 'green');
            flagID = true;
        }
    }

    // 비밀번호 유효성 검사
    function validateInputPw() {
        const regexPw = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,16}$/;
        checkPW = $('#password').val();

        if (checkPW === "") {
            $('#pw-message').html('<i class="fa-solid fa-x"></i> 비밀번호를 입력해주세요.');
            $('#pw-message').css('color', 'red');
            flagPW = false;
            return;
        }

        if (!regexPw.test(checkPW)) {
            $('#pw-message').html('<i class="fa-solid fa-x"></i> 비밀번호는 8~16자의 영문 대/소문자, 숫자, 특수문자(!,@,#,$,%,^,&,*,-)를 사용해주세요.');
            $('#pw-message').css('color', 'red');
            flagPW = false;
        } else {
            $('#pw-message').html('<i class="fa-solid fa-check"></i> 사용 가능한 비밀번호 형식입니다.');
            $('#pw-message').css('color', 'green');
            flagPW = true;
        }
    }

    // 비밀번호 보기/숨기기
    $('#togglePassword').on('click', function () {
        const passwordField = $('#password');
        const type = passwordField.attr('type') === 'password' ? 'text' : 'password';
        passwordField.attr('type', type);
        $(this).text(type === 'password' ? '보기' : '숨기기');
    });

    // 이름 유효성 검사
    function validateInputName() {
        const regexName = /^[가-힣]{2,10}$/;
        checkName = $('#user_name').val();

        if (checkName === "") {
            $('#name-message').html('<i class="fa-solid fa-x"></i> 이름을 입력해주세요.');
            $('#name-message').css('color', 'red');
            flagName = false;
            return;
        }

        if (!regexName.test(checkName)) {
            $('#name-message').html('<i class="fa-solid fa-x"></i> 이름은 2~10자의 한글만 사용 가능합니다.');
            $('#name-message').css('color', 'red');
            flagName = false;
        } else {
            $('#name-message').html('<i class="fa-solid fa-check"></i> 사용 가능한 이름 형식입니다.');
            $('#name-message').css('color', 'green');
            flagName = true;
        }
    }

    // 이메일 유효성 검사
    function validateInputEmail() {
        const regexEmail = /(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))/
        checkEmail = $('#email').val();

        if (checkEmail === "") {
            $('#email-message').html('<i class="fa-solid fa-x"></i> 이메일을 입력해주세요.');
            $('#email-message').css('color', 'red');
            flagEmail = false;
            return;
        }

        if (!regexEmail.test(checkEmail)) {
            $('#email-message').html('<i class="fa-solid fa-x"></i> 이메일 형식이 올바르지 않습니다.');
            $('#email-message').css('color', 'red');
            flagEmail = false;
        } else {
            $('#email-message').html('<i class="fa-solid fa-check"></i> 사용 가능한 이메일 형식입니다.');
            $('#email-message').css('color', 'green');
            flagEmail = true;
        }
    }

    // 전화번호 유효성 검사
    function validateInputPhone() {
        const regexPhone = /^\d{3}-\d{3,4}-\d{4}$/;
        checkPhone = $('#phone_number').val();

        if (checkPhone === "") {
            $('#phone-message').html('<i class="fa-solid fa-x"></i> 전화번호를 입력해주세요.');
            $('#phone-message').css('color', 'red');
            flagPhone = false;
            return;
        }

        if (!regexPhone.test(checkPhone)) {
            $('#phone-message').html('<i class="fa-solid fa-x"></i> 전화번호 형식이 올바르지 않습니다.');
            $('#phone-message').css('color', 'red');
            flagPhone = false;
        } else {
            $('#phone-message').html('<i class="fa-solid fa-check"></i> 사용 가능한 전화번호 형식입니다.');
            $('#phone-message').css('color', 'green');
            flagPhone = true;
        }
    }


    // 아이디 중복 체크
    $('#checkID').on('click', function (e) {
        e.preventDefault();
        e.stopPropagation();  // 이벤트 전파 중지

        const userId = $('#user_id').val();
        console.log("아이디", userId);
        if (flagID === true) {
            $.ajax({
                url: '/join/check-duplicate',
                method: 'POST',
                contentType: 'application/json; charset=utf-8',
                dataType: 'json',
                data: JSON.stringify({
                    USER_ID: userId
                })
            })
                .done(function (response) {
                    console.log("서버 응답:", response);
                    if (response.exists) {
                        $('#id-message').html('<i class="fa-solid fa-x"></i> 이미 사용 중인 아이디입니다.');
                        $('#id-message').css('color', 'red');
                        alert('이미 사용 중인 아이디입니다.');
                    } else {
                        alert('사용 가능한 아이디입니다.');
                    }
                })
                .fail(function (jqXHR, textStatus, errorThrown) {
                    console.error("AJAX 오류:", textStatus, errorThrown);
                    alert('중복 확인 중 오류가 발생했습니다.');
                });
        } else {
            alert('올바른 아이디 형식을 입력해주세요.');
        }
    });

    // 회원가입 폼 제출
    $('#joinBtn').on('click', function (e) {
        e.preventDefault();
        console.log(flagID);
        if (flagID === true && flagPW === true && flagName === true && flagEmail === true && flagPhone === true) {
            $.ajax({
                url: '/join/submit',
                method: 'POST',
                contentType: 'application/json; charset=utf-8',
                dataType: 'json',
                data: JSON.stringify({
                    USER_ID: $('#user_id').val(),
                    PASSWORD: $('#password').val(),
                    USER_NAME: $('#user_name').val(),
                    EMAIL: $('#email').val(),
                    PHONE_NUMBER: $('#phone_number').val()
                })
            })
                .done(function (response) {
                    console.log("서버 응답 :", response);
                    alert('회원가입이 완료되었습니다.');
                    window.location.href = '/';
                })
                .fail(function (jqXHR, textStatus, errorThrown) {
                    console.error("AJAX 오류:", textStatus, errorThrown);
                    if (jqXHR.status === 409) {
                        const errorMessage = jqXHR.responseJSON.error;
                    } else {
                        alert('회원가입 중 오류가 발생했습니다. 다시 시도해주세요.');
                    }
                });
        } else {
            alert('모든 필드를 올바르게 입력해주세요.');
        }
    });
});