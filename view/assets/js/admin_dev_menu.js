// DOM 요소 선언
let input;
let preview;

// DOM이 완전히 로드된 후 실행
$(document).ready(function () {
    // DOM 요소 초기화
    input = document.getElementById('image_uploads');
    preview = document.getElementById('dev-img');

    // 요소가 존재하는지 확인 후 이벤트 리스너 등록
    if (input && preview) {
        input.addEventListener('change', function (e) {
            const file = e.target.files[0];
            if (file && validFileType(file)) {
                const reader = new FileReader();
                reader.onload = function (e) {
                    preview.src = e.target.result;
                }
                reader.readAsDataURL(file);
                console.log('이미지 프리뷰 업데이트 완료!!!');
            } else {
                alert('이미지를 불러올 수 없습니다!!!');
            }
        });
    } else {
        console.error('이미지 업로드 요소를 찾을 수 없습니다');
    }
});

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
    const dev_main_exposure = $('#dev_main_exposure').val() === 'true' ? true : false
    formData.append('DEV_NAME', $('#dev_name').val());
    formData.append('DEV_DETAIL_NAME', $('#dev_detail_name').val());
    formData.append('DEV_HISTORY', $('#dev_history').val());
    formData.append('VIEW_YN', dev_main_exposure);

    // FormData 내용 확인
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
            alert("개발자 등록에 성공했습니다");
            $('#dev-add').modal('hide');
            location.reload();
        },
        error: function (error) {
            console.error('에러:', error);
            alert("개발자 등록 실패");
        }
    });
}