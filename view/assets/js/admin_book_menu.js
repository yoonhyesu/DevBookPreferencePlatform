import { Tabulator } from '/assets/js/admin_grid.js';
import { selectedData } from '/assets/js/admin_grid.js';

// 도서 관리페이지
$(document).ready(function () {

    // 도서 태그 값 불러오기
    let multipleSelect = new Choices('#bookTag', {
        removeItemButton: true,
        searchEnabled: true,
        placeholder: true,
        placeholderValue: '태그를 선택해주세요(최대 5개)'
    });

    $.ajax({
        url: '/admin/book/tags',
        method: 'GET',
        success: function (tags) {
            const choices = tags.map(tag => ({
                value: tag.TAG_ID.toString(),
                label: tag.TAG_NAME
            }));
            multipleSelect.setChoices(choices, 'value', 'label', true);
        },
        error: function (xhr, status, error) {
            console.error('태그 조회 실패:', error);
            alert('태그 조회 실패');
        }
    });

    // 도서 검색 관련 코드
    $('#searchBookBtn').click(function () {
        const searchTitle = $('#title').val().trim();
        if (!searchTitle) {
            alert('검색할 도서명을 입력해주세요.');
            return;
        }

        // 현재 활성화된 요소의 포커스 해제
        $(document.activeElement).blur();

        // API 호출 URL 수정
        $.ajax({
            url: `/admin/book/search?title=${encodeURIComponent(searchTitle)}`,
            method: 'GET',
            success: function (response) {
                console.log('API 응답:', response);

                if (response && response.docs && response.docs.length > 0) {
                    const bookAddModal = bootstrap.Modal.getInstance($('#book-add'));
                    if (bookAddModal) {
                        bookAddModal.hide();
                    }

                    displaySearchResults(response.docs);
                } else {
                    alert('검색 결과가 없습니다');
                    $('#book-search').modal('hide');
                    $('#book-add').modal('show');
                }
            },
            error: function (xhr, status, error) {
                console.error('검색 오류:', error);
                alert('검색 중 오류가 발생했습니다');
            }
        });
    });

    // 검색 결과 표시 함수
    function displaySearchResults(books) {
        const searchResultsContainer = $('#book-search .modal-body .row');
        searchResultsContainer.empty();

        books.forEach((book, index) => {
            // book 데이터를 Base64로 인코딩
            const encodedBookData = btoa(encodeURIComponent(JSON.stringify(book)));

            const bookHtml = `
            <div class="col-12 col-sm-6 col-lg-4 d-flex align-items-center justify-content-center flex-column">
                <img src="${book.TITLE_URL || '/assets/images/no-image.jpg'}" 
                     width="200" height="250" 
                     class="bd-placeholder-img border" 
                     alt="${book.TITLE}">
                <div class="mt-2">
                    <input class="form-check-input" 
                           type="radio" 
                           name="bookSelect" 
                           value="${index}" 
                           data-book="${encodedBookData}">
                </div>
            </div>
        `;
            searchResultsContainer.append(bookHtml);
        });

        $('#book-search').modal('show');
    }

    // 도서 선택시 이벤트
    $('#selectBookBtn').click(function () {
        const selectedRadio = $('input[name="bookSelect"]:checked');
        if (selectedRadio.length === 0) {
            alert('선택된 도서가 없습니다');
            return;
        }

        // Base64 디코딩 후 JSON 파싱(내용에 한글이나 특수문자가 포함될 수 있음)
        try {
            const encodedData = selectedRadio.attr('data-book');
            const decodedData = decodeURIComponent(atob(encodedData));
            const selectedBook = JSON.parse(decodedData);
            console.log("선택한 도서 정보", selectedBook);

            // book-add 모달의 입력 필드에 데이터 바인딩
            $('#title').val(selectedBook.TITLE);
            $('#book_title').val(selectedBook.TITLE);
            $('#author').val(selectedBook.AUTHOR);
            $('#isbn').val(selectedBook.EA_ISBN);
            $('#isbnAddCode').val(selectedBook.EA_ADD_CODE);
            $('#publisher').val(selectedBook.PUBLISHER);
            $('#publishPredate').val(selectedBook.PUBLISH_PREDATE);
            $('#imageUrl').val(selectedBook.TITLE_URL);
            const PRE_PRICE_VAL = selectedBook.PRE_PRICE.replace(/[^\d]/g, '');
            $('#prePrice').val(PRE_PRICE_VAL);
            $('#ebookYn').val(selectedBook.EBOOK_YN);
            $('#titleUrl').val(selectedBook.TITLE_URL);
            $('#bookTbCntUrl').val(selectedBook.BOOK_TB_CNT);
            $('#bookIntroductionUrl').val(selectedBook.BOOK_INTRODUCTION);
            $('#bookSummaryUrl').val(selectedBook.BOOK_SUMMARY);
            $('#page').val(selectedBook.PAGE);

            // 모달 전환
            $('#book-search').modal('hide');
            $('#book-add').modal('show');
        } catch (error) {
            console.error('도서 데이터 파싱 오류:', error);
            alert('도서 데이터 처리 중 오류가 발생했습니다');
        }
    });

    // 선택된 도서 데이터 베이스 전송
    $('#sendBookData').click(function () {
        // 문자열 정리 함수
        const cleanString = (str) => str ? str.trim().replace(/\s+/g, ' ') : '';

        // 필수 필드 검증
        const requiredFields = ['book_title', 'author', 'isbn'];
        const emptyFields = requiredFields.filter(fieldId => !$(`#${fieldId}`).val().trim());

        if (emptyFields.length > 0) {
            alert('필수 정보를 모두 입력해주세요');
            return;
        }

        const devContents = [];

        // 첫 번째 추천 프로그래머 정보가 있는 경우에만 추가
        const firstDevId = $('#devID1').val();
        const firstReason = $('#reason1').val();
        if (firstDevId && firstReason) {
            devContents.push({
                DEV_ID: firstDevId,
                DEV_RECOMMEND_REASON: firstReason,
            });
        }

        // 추가된 추천 프로그래머 정보 중 유효한 데이터만 추가
        const containers = $('#recommenderContainer').children();
        for (let i = 0; i < containers.length / 3; i++) {
            const devId = $(`#devID${i + 2}`).val();
            const reason = $(`#reason${i + 2}`).val();

            if (devId && reason) { // ID와 추천이유가 모두 있는 경우에만 추가
                devContents.push({
                    DEV_ID: devId,
                    DEV_RECOMMEND_REASON: reason,
                });
            }
        }

        // 최소 1명 이상의 추천 프로그래머 필요
        if (devContents.length === 0) {
            alert('최소 1명 이상의 추천 프로그래머가 필요합니다');
            return;
        }

        // 데이터 형식 변환
        const sendData = {
            book: {
                BOOK_TITLE: cleanString($('#book_title').val()),
                PUBLISHER: cleanString($('#publisher').val()),
                AUTHOR: cleanString($('#author').val()),
                PUBLISH_DATE: $('#publishPredate').val(),
                PRICE: parseInt($('#prePrice').val()) || 0,
                CONTENTS_LIST: cleanString($('#bookTbCntUrl').val()),
                COVER_URL: cleanString($('#titleUrl').val()),
                SUMMARY: cleanString($('#bookSummaryUrl').val()),
                DESCRIPTION: cleanString($('#bookIntroductionUrl').val()),
                RECOMMENDATION: cleanString($('#recommendation').val()),
                TAGS: $('#bookTag').val().length ? '&' + $('#bookTag').val().join('&') + '&' : '',
                GRADE: cleanString($('#grade').val()),
                ISBN: cleanString($('#isbn').val()),
                DEL_YN: false,
                ISBN_ADD: cleanString($('#isbnAddCode').val()),
                PAGE: parseInt($('#page').val()) || 0,
                EBOOK_YN: cleanString($('#ebookYn').val())
            },
            DEV_CONTENTS: devContents
        };

        // AJAX 요청
        $.ajax({
            url: '/admin/book/add',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(sendData),
            success: function (response) {
                alert('책 데이터가 성공적으로 저장되었습니다');
                $('#book-add').modal('hide');
                location.reload();
            },
            error: function (xhr, status, error) {
                console.error('전송된 데이터:', sendData);
                console.error('에러:', error);
                alert('책 데이터 저장 중 오류가 발생했습니다');
            }
        });
    });


    // 모달 닫기 버튼 이벤트
    $('.btn-close, .modal .btn-secondary').click(function () {
        const modal = $(this).closest('.modal');
        modal.modal('hide');

        // 모든 입력값 초기화
        modal.find('input').val('');

        // 검색 결과 컨테이너 비우기
        if (modal.attr('id') === 'book-search') {
            $('#book-search .modal-body .row').empty();
        }

        // 라디오 버튼 선택 해제
        $('input[name="bookSelect"]').prop('checked', false);

        // 새로고침
        location.reload();
    });

    // 추천 프로그래머 추가 버튼 이벤트
    $('#add-recommender').click(function () {
        const index = ($('#recommenderContainer').children().length / 3) + 1;
        $('#recommenderContainer').append(`
            <div class="col-12 mb-3">
                <label for="devID${index}" class="form-label required">추천 프로그래머ID</label>
                <div class="input-group">
                    <input type="text" class="form-control" id="devID${index}" name="devID${index}" readonly required>
                    <button class="btn btn-outline-primary checkDevId" type="button" id="checkDevId${index}"
                        data-toggle="modal" data-target="#dev-search-modal">검색</button>
                </div>
            </div>
            <div class="col-12 mb-3">
                <label for="name${index}" class="form-label required">추천 프로그래머명</label>
                <input type="text" class="form-control" id="name${index}" name="name${index}" readonly required>
            </div>
            <div class="col-12 mb-3">
                <label for="reason${index}" class="form-label required">프로그래머 추천이유</label>
                <textarea class="form-control" id="reason${index}" name="reason${index}" rows="3" required></textarea>
            </div>
        `);

        // 새로 추가된 검색 버튼에 이벤트 핸들러 연결
        $(`#checkDevId${index}`).on("click", function () {
            const currentIndex = $(this).attr('id').replace('checkDevId', '');
            $('#dev-search-modal').data('currentIndex', currentIndex);
            const searchModal = new bootstrap.Modal(document.getElementById('dev-search-modal'));
            searchModal.show();
        });
    });

    // 추천 프로그래머 삭제 버튼 이벤트
    $('#del-recommender').click(function () {
        $('#recommenderContainer').children().last().remove();
        $('#recommenderContainer').children().last().remove();
        $('#recommenderContainer').children().last().remove();
    });

    // 프로그래머 검색 버튼 클릭 이벤트
    $('#checkDevId').on("click", function () {
        // 클릭된 검색 버튼의 가장 가까운 입력 필드의 index 저장
        const $inputGroup = $(this).closest('.input-group');
        const $devIdInput = $inputGroup.find('input[id^="devID"]');
        const currentIndex = $devIdInput.attr('id').replace('devID', '');

        // 현재 index를 modal에 데이터 속성으로 저장
        $('#dev-search-modal').data('currentIndex', currentIndex);

        // Bootstrap 5 방식으로 모달 표시
        const searchModal = new bootstrap.Modal(document.getElementById('dev-search-modal'));
        searchModal.show();
    });

    // 프로그래머 선택 버튼 클릭 이벤트
    $('#selectDevBtn').on("click", function () {
        if (!selectedData) {
            alert('프로그래머를 선택해주세요');
            return;
        }

        // 모달에 저장된 현재 index 가져오기
        const currentIndex = $('#dev-search-modal').data('currentIndex');

        // 선택된 프로그래머 정보를 해당 index의 입력 필드에 설정
        $(`#devID${currentIndex}`).val(selectedData.ID);
        $(`input[name="name${currentIndex}"]`).val(selectedData.DEV_NAME);

        // 모달 닫기
        const searchModal = bootstrap.Modal.getInstance(document.getElementById('dev-search-modal'));
        searchModal.hide();
    });

    // 페이지 로드 시 첫 번째 검색 버튼에 이벤트 핸들러 연결
    $(document).ready(function () {
        $('#checkDevId1').on("click", function () {
            $('#dev-search-modal').data('currentIndex', '1');
            const searchModal = new bootstrap.Modal(document.getElementById('dev-search-modal'));
            searchModal.show();
        });
    });
});



