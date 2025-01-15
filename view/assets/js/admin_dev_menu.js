$(document).ready(function () {



})

// 개발자 등록 버튼 이벤트
function dev_add() {
    const dev_name = $('#dev_name').val();
    const dev_nickname = $('#dev_detail_name').val();
    const dev_history = $('#dev_history').val();
    const dev_main_exposure = $('#dev_main_exposure').val() === "true" ? true : false

    $.ajax({
        url: '/admin/dev/add',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            DEV_NAME: dev_name,
            DEV_DETAIL_NAME: dev_nickname,
            DEV_HISTORY: dev_history,
            VIEW_YN: dev_main_exposure
        }),
        success: function (response) {
            alert("개발자 등록에 성공했습니다");
            $('#dev-add').modal('hide');
            location.reload();
        },
        error: function (error) {
            alert("개발자 등록 실패");
        }
    })
}