var input = document.querySelector('#image_uploads');
var preview = document.getElementById('dev-img');

$(function () {
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
    'image/jpg'
]

// 이미지 파일 유효성 검사
function validFileType(file) {
    return fileTypes.includes(file.type);
}

// 개발자 등록 버튼 이벤트
function dev_add() {
    const formData = new FormData();
    // 이미지 파일 있는지 확인
    const imageFile = input.files[0];
    if (imageFile) {
        console.log("이미지 파일 첨부:", imageFile);
        formData.append('PROFILE_IMAGE', imageFile);
    }

    // 나머지 데이터 추가
    const dev_main_exposure = $('#dev_main_exposure').val() === "true" ? true : false
    formData.append('DEV_NAME', $('#dev_name').val());
    formData.append('DEV_DETAIL_NAME', $('#dev_detail_name').val());
    formData.append('DEV_HISTORY', $('#dev_history').val());
    formData.append('VIEW_YN', dev_main_exposure);

    for (let pair of formData.entries()) {
        console.log(pair[0] + ': ' + pair[1]);
    }

    $.ajax({
        url: '/admin/dev/add',
        method: 'POST',
        processData: false,
        contentType: false,
        data: formData,
        success: function (response) {
            console.log('서버응답;', response);
            if (response.PROFILE_IMAGE_PATH) {
                $('#profile-img').attr('src', '/storage' + response.PROFILE_IMAGE_PATH);
            }
            alert("개발자 등록에 성공했습니다");
            $('#dev-add').modal('hide');
            location.reload();
        },
        error: function (error) {
            alert("개발자 등록 실패");
        }
    })
}