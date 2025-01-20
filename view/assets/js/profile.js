var input = document.querySelector('#image_uploads');
var preview = document.getElementById('profile-img');

$(function () {
    UserInfoLoad();
    input.addEventListener('change', updateImageDisplay);
})

// 프로필 업로드 미리보기
function updateImageDisplay() {
    const file = input.files[0];
    if (file && validFileType(file)) {
        preview.src = URL.createObjectURL(file);
    } else {
        alert('이미지를 불러올 수 없습니다');
    }
}

const fileTypes = [
    'image/jpeg',
    'image/png',
]
// 이미지 파일 유효성 검사
function validFileType(file) {
    return fileTypes.includes(file.type);
}

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

        // 이미지 경로가 있을 때만 설정하도록 수정
        if (response.PROFILE_IMAGE_PATH) {
            const imagePath = '/storage/image/dev' + response.PROFILE_IMAGE_PATH.replace(/\\/g, '/');
            preview.src = imagePath;
            console.log('변환된 이미지 경로:', imagePath);
        } else {
            preview.src = '/assets/images/profile.png';
        }
    }).fail(function (jqXHR, textStatus, errorThrown) {
        console.error('Ajax 요청 실패:', textStatus, errorThrown); // 에러 로깅 추가!!!
    });
}

//프로필 수정
function updateProfile() {
    const formData = new FormData();

    // 이미지 파일이 있는지 확인
    const imageFile = input.files[0];
    console.log("이미지파일", imageFile)
    if (imageFile) {
        console.log("이미지 파일 첨부:", imageFile); // 디버깅용
        formData.append('PROFILE_IMAGE', imageFile);
    }

    // 나머지 데이터 추가
    formData.append('USER_NAME', $('#user_name').val());
    formData.append('PHONE_NUMBER', $('#userPhone').val());
    formData.append('USER_TEXT', $('#user_text').val());
    formData.append('COMPANY', $('#user_company').val());
    formData.append('GITHUB_LINK', $('#github_url').val());
    formData.append('BLOG_LINK', $('#blog_url').val());

    // FormData 내용 확인
    for (let pair of formData.entries()) {
        console.log(pair[0] + ': ' + pair[1]); // 디버깅용
    }

    $.ajax({
        url: '/account/profile-edit/manage',
        method: 'POST',
        processData: false, // 중요!!!
        contentType: false, // 중요!!!
        data: formData,
        success: function (response) {
            console.log("서버 응답:", response); // 디버깅용
            if (response.PROFILE_IMAGE_PATH) {
                $('#profile-img').attr('src', '/storage' + response.PROFILE_IMAGE_PATH);
            }
            alert('프로필이 성공적으로 수정되었습니다');
            location.reload();
        },
        error: function (error) {
            console.error('프로필 수정 실패:', error);
            alert('프로필 수정에 실패했습니다');
        }
    });
}

// 수정 버튼 클릭 이벤트 핸들러
$(document).ready(function () {
    $('#profile_update_btn').on('click', function (e) {
        e.preventDefault();
        updateProfile();
    });
});

