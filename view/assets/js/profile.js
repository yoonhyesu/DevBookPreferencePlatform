$(document).ready(function () {
    UserInfoLoad();

});

function UserInfoLoad() {
    $.ajax({
        url: '/account/profile-edit',
        method: 'GET',
        dataType: 'json'
    }).done(function (response) {
        console.log('프로필 조회 응답:', response);
        if (response.error) {
            alert(response.error);
            return;
        }
        $('#userPhone').val(response.PHONE_NUMBER || '');
        $('#user_text').val(response.USER_TEXT || '');
        $('#user_user_company').val(response.COMPANY || '');
        $('#github_url').val(response.GITHUB_LINK || '');
        $('#blog_url').val(response.BLOG_LINK || '');
        $('#user_company').val(response.COMPANY || '');
    }).fail(function (error) {
        console.error('프로필 조회 실패:', error);
        alert('프로필 조회에 실패했습니다!!!');
    });
}

// 프로필 수정 함수
function updateProfile() {
    const profileData = {
        USER_NAME: $('#user_name').val(),
        PHONE_NUMBER: $('#userPhone').val(),
        USER_TEXT: $('#user_text').val(),
        COMPANY: $('#user_company').val(),
        GITHUB_LINK: $('#github_url').val(),
        BLOG_LINK: $('#blog_url').val()
    };

    $.ajax({
        url: '/account/profile-edit/manage',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(profileData),
        dataType: 'json'
    }).done(function (response) {
        alert('프로필이 성공적으로 수정되었습니다!!!');
        location.reload();
    }).fail(function (error) {
        console.error('프로필 수정 실패:', error);
        alert('프로필 수정에 실패했습니다!!!');
    });
}

// 수정 버튼 클릭 이벤트 핸들러
$(document).ready(function () {
    $('#profile_update_btn').on('click', function (e) {
        e.preventDefault();
        updateProfile();
    });
});