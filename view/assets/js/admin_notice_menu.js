import { Tabulator } from '/assets/js/admin_grid.js';
import { selectedData } from '/assets/js/admin_grid.js';

// 공지 관리페이지

// 전역 스코프에 함수 등록
window.notice_add = function () {
    const title = $('#notice_title').val();
    const content = $('#notice_content').val();
    const topyn = $('#notice_topyn').val();

    if (!title || !content) {
        alert('제목과 내용을 모두 입력해주세요');
        return;
    }

    $.ajax({
        url: '/admin/notice/add',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            TITLE: title,
            CONTENT: content,
            TOP_YN: topyn === 'true'
        }),
        success: function (response) {
            alert("공지 등록에 성공했습니다");
            $('#notice-add').modal('hide');
            location.reload();
        },
        error: function (error) {
            alert('공지 등록에 실패했습니다');
            console.error('에러:', error);
        }
    });
};

// 입력값 초기화 함수도 전역으로
window.clearNoticeInputs = function () {
    $('#notice_title').val('');
    $('#notice_content').val('');
    $('#notice_topyn').val('false');
};

// DOM이 준비되면 실행
$(document).ready(function () {
    // 모달이 닫힐 때 입력값 초기화
    $('#notice-add').on('hidden.bs.modal', function () {
        clearNoticeInputs();
    });
});