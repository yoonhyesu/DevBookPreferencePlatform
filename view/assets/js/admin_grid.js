import { TabulatorFull as Tabulator } from '/node_modules/tabulator-tables/dist/js/tabulator_esm.min.js';

export { Tabulator };
export let selectedData = null;
export let admin_dev_search_table = null;

// 검색 관련 함수들을 모듈 스코프 밖으로 export
export function allSearch(data, searchWord) {
    if (Object.keys(searchWord).Length == 0) {
        return true;
    }
    var devName = data.DEV_NAME.toLowerCase().indexOf(searchWord.toLowerCase());
    var devID = data.ID.toLowerCase().indexOf(searchWord.toLowerCase());
    return (devName != -1 ? true : false) || (devID != -1 ? true : false);
}

export function updateFilter() {
    if (!admin_dev_search_table) return;

    var fieldValue = $('#filter-field').val();
    var searchWord = $('#filter-value').val();

    if (fieldValue == 'ALL') {
        admin_dev_search_table.setFilter(allSearch, searchWord);
    }
    else if (fieldValue) {
        admin_dev_search_table.setFilter(fieldValue, 'like', searchWord);
    }
}

export function initSearch() {
    if (!admin_dev_search_table) return;

    $('#filter-field').val('ALL');
    $('#filter-value').val('');
    admin_dev_search_table.clearFilter();
}

$(document).ready(function () {
    let admin_notice_table, admin_book_table, admin_dev_table;

    // 공지사항 관리
    if (document.getElementById("admin_notice_table")) {
        admin_notice_table = new Tabulator("#admin_notice_table", {
            height: "1000px",
            layout: "fitColumns",
            selectableRows: true,
            rowHeader: {
                headerSort: false, resizable: false, frozen: true, headerHozAlign: "center", hozAlign: "center", formatter: "rowSelection", titleFormatter: "rowSelection", cellClick: function (e, row) {
                    var rowData = row.getData();
                },
            },
            placeholder: "데이터없음",
            ajaxURL: "/admin/notice/manage",  // URL 수정
            ajaxConfig: {
                method: "GET",
                headers: {
                    "Content-type": 'application/json; charset=utf-8',
                },
                credentials: "include"
            },
            pagination: true,
            paginationSize: 10,
            columns: [
                { title: "no.", field: "NOTICE_ID", sorter: "number", width: 100 },
                { title: "제목", field: "TITLE", sorter: "string", width: 300 },
                { title: "내용", field: "CONTENT", sorter: "string", width: 300 },
                { title: "상위여부", field: "TOP_YN", sorter: "string", width: 150, formatter: "tickCross" },
                {
                    title: "등록일",
                    field: "CREATED_DATE",
                    sorter: "datetime",
                    width: 200,
                    formatter: function (cell) {
                        const date = new Date(cell.getValue());
                        return date.toLocaleDateString('ko-KR', {
                            year: 'numeric',
                            month: '2-digit',
                            day: '2-digit',
                            hour: '2-digit',
                            minute: '2-digit'
                        });
                    }
                },
                {
                    title: "수정일",
                    field: "UPDATED_DATE",
                    sorter: "datetime",
                    width: 200,
                    formatter: function (cell) {
                        const date = new Date(cell.getValue());
                        return date.toLocaleDateString('ko-KR', {
                            year: 'numeric',
                            month: '2-digit',
                            day: '2-digit',
                            hour: '2-digit',
                            minute: '2-digit'
                        });
                    }
                }
            ],
        });
        admin_notice_table.on("ajaxError", function (error) {
            console.error("Ajax 로딩 에러:", error);
        });

        admin_notice_table.on("rowSelected", function (row) {
            console.log("선택된 행 데이터:", row.getData());
            selectedData = row.getData();
        });
    }


    // 추천 도서 관리
    if (document.getElementById("admin_book_table")) {
        admin_book_table = new Tabulator("#admin_book_table", {
            height: "1000px",
            width: "50%",
            layout: "fitColumns",
            selectableRows: true,
            rowHeader: {
                headerSort: false, resizable: false, frozen: true, headerHozAlign: "center", hozAlign: "center", formatter: "rowSelection", titleFormatter: "rowSelection", cellClick: function (e, row) {
                    var rowData = row.getData();
                },
            },
            placeholder: "데이터없음",
            ajaxURL: "/admin/book/manage",
            ajaxConfig: {
                method: "GET",
                headers: {
                    "Content-type": 'application/json; charset=utf-8',
                    "X-CSRF-TOKEN": $('meta[name="csrf-token"]').attr('content')
                },
                credentials: "include"
            },
            pagination: true,
            paginationMode: "remote",
            paginationSize: 10,
            paginationInitialPage: 1,
            ajaxResponse: function (url, params, response) {
                console.log("Server Response:", response);
                return {
                    data: response,
                    last_page: 1,
                    total: response.length
                };
            },
            columns: [
                { title: "no.", field: "BOOK_ID", sorter: "number", width: 100 },
                { title: "도서명", field: "BOOK_TITLE", sorter: "string", width: 400 },
                { title: "저자", field: "AUTHOR", sorter: "string", width: 300 },
                { title: "ISBN", field: "ISBN", sorter: "string", width: 150 },
                { title: "발매일", field: "PUBLISH_DATE", sorter: "string", width: 200 },
                {
                    title: "표지",
                    field: "COVER_URL",
                    sorter: "string",
                    formatter: "image",
                    formatterParams: {
                        height: "100px",
                        width: "80px"
                    },
                    width: 100
                }
            ],
        });

        admin_book_table.on("ajaxError", function (error) {
            console.error("Ajax 로딩 에러:", error);
        });

        admin_book_table.on("rowClick", function (e, row) {
            console.log("클릭된 행 데이터:", row.getData());
            selectedData = row.getData();

        });
    }

    // 개발자 관리 
    if (document.getElementById("admin_dev_table")) {
        admin_dev_table = new Tabulator("#admin_dev_table", {
            height: "1000px",
            width: "50%",
            layout: "fitColumns",
            selectableRows: true,
            rowHeader: {
                headerSort: false, resizable: false, frozen: true, headerHozAlign: "center", hozAlign: "center", formatter: "rowSelection", titleFormatter: "rowSelection", cellClick: function (e, row) {
                    var rowData = row.getData();
                },
            },
            placeholder: "데이터없음",
            ajaxURL: "/admin/dev/manage",
            ajaxConfig: {
                method: "GET",
                headers: {
                    "Content-type": 'application/json; charset=utf-8',
                    "X-CSRF-TOKEN": $('meta[name="csrf-token"]').attr('content')
                },
                credentials: "include"
            },
            pagination: true,
            paginationMode: "remote",
            paginationSize: 10,
            paginationInitialPage: 1,
            ajaxResponse: function (url, params, response) {
                console.log("Server Response:", response);
                return {
                    data: response,
                    last_page: 1,
                    total: response.length
                };
            },
            columns: [
                { title: "no.", field: "ID", sorter: "string", width: 100 },
                { title: "프로그래머명", field: "DEV_NAME", sorter: "string", width: 300 },
                { title: "별칭", field: "DEV_DETAIL_NAME", sorter: "string", width: 300 },
                { title: "경력", field: "DEV_HISTORY", sorter: "string", width: 400 },
                { title: "메인노출여부", field: "VIEW_YN", sorter: "string", width: 200, formatter: "tickCross" },
            ],
        });

        admin_dev_table.on("ajaxError", function (error) {
            console.error("Ajax 로딩 에러:", error);
        });

        admin_dev_table.on("rowClick", function (e, row) {
            console.log("클릭된 행 데이터:", row.getData());
            selectedData = row.getData();
        });
    }

    // 개발자 검색
    if (document.getElementById("admin_dev_search_table")) {
        admin_dev_search_table = new Tabulator("#admin_dev_search_table", {
            height: "1000px",
            width: "50%",
            layout: "fitColumns",
            selectableRows: true,
            rowHeader: {
                headerSort: false, resizable: false, frozen: true, headerHozAlign: "center", hozAlign: "center", formatter: "rowSelection", titleFormatter: "rowSelection", cellClick: function (e, row) {
                    var rowData = row.getData();
                },
            },
            placeholder: "데이터없음",
            ajaxURL: "/admin/dev/manage",
            ajaxConfig: {
                method: "GET",
                headers: {
                    "Content-type": 'application/json; charset=utf-8',
                    "X-CSRF-TOKEN": $('meta[name="csrf-token"]').attr('content')
                },
                credentials: "include"
            },
            pagination: true,
            paginationMode: "remote",
            paginationSize: 10,
            paginationInitialPage: 1,
            ajaxResponse: function (url, params, response) {
                console.log("Server Response:", response);
                return {
                    data: response,
                    last_page: 1,
                    total: response.length
                };
            },
            columns: [
                { title: "no.", field: "ID", sorter: "string", width: 300 },
                { title: "프로그래머명", field: "DEV_NAME", sorter: "string", width: 200 },
                { title: "메인노출여부", field: "VIEW_YN", sorter: "string", width: 200, formatter: "tickCross" },
            ],
        });

        admin_dev_search_table.on("ajaxError", function (error) {
            console.error("Ajax 로딩 에러:", error);
        });

        admin_dev_search_table.on("rowClick", function (e, row) {
            console.log("클릭된 행 데이터:", row.getData());
            selectedData = row.getData();
        });
    }

    // 개발자 관리 - 수정 버튼 클릭시 모달에 바인딩
    $('#update_DEV').on('click', function (e) {
        console.log("현재 selectedData:", selectedData);
        if (!selectedData) {
            e.preventDefault();
            alert('수정할 개발자를 선택해주세요.');
            return;
        }

        // 데이터 바인딩
        $('#u_dev_name').val(selectedData.DEV_NAME);
        $('#u_dev_detail_name').val(selectedData.DEV_DETAIL_NAME);
        $('#u_dev_history').val(selectedData.DEV_HISTORY);
        $('#u_dev_main_exposure').val(selectedData.VIEW_YN === true ? 'true' : 'false');

        // Bootstrap 5 방식으로 모달 표시
        const updateModal = new bootstrap.Modal(document.getElementById('dev-update'));
        updateModal.show();
    });

    // 개발자 관리 - 수정 모달 서버에 전송
    $('#update-btn').on('click', function (e) {
        const dev_name = $('#u_dev_name').val();
        const dev_detail_name = $('#u_dev_detail_name').val();
        const dev_history = $('#u_dev_history').val();
        const dev_main_exposure = $('#u_dev_main_exposure').val();

        $.ajax({
            url: '/admin/dev/update',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                NOTICE_ID: selectedData.NOTICE_ID,
                DEV_NAME: dev_name,
                DEV_DETAIL_NAME: dev_detail_name,
                DEV_HISTORY: dev_history,
                VIEW_YN: dev_main_exposure === 'true' ? true : false
            }),
            success: function (response) {
                alert("개발자 수정에 성공했습니다");
                $('#dev-update').modal('hide');
            },
            error: function (error) {
                alert("개발자 수정 실패");
            }
        });
    });

    // 개발자 관리 - 모달이 닫힐 때 초기화
    $('#dev-update').on('hidden.bs.modal', function (e) {
        // 모달 내의 모든 input과 textarea 초기화
        $(this).find('input, textarea, select').val('');
        // 선택된 데이터도 초기화
        selectedData = null;
        // 사이트 리로드
        location.reload();
    });

    // 공지사항 수정 
    $('#update_notice').on('click', function (e) {
        if (!selectedData) {
            e.preventDefault();
            alert('수정할 공지사항을 선택해주세요!!!');
            return;
        }

        // 데이터 바인딩 - HTML의 ID와 맞춤
        $('#notice_update_title').val(selectedData.TITLE);
        $('#notice_update_content').val(selectedData.CONTENT);
        $('#notice_update_topyn').val(selectedData.TOP_YN === true ? 'true' : 'false');

        // Bootstrap 5 모달 표시
        const updateModal = new bootstrap.Modal(document.getElementById('notice-update'));
        updateModal.show();
    });

    // 공지사항 수정 실행
    window.notice_update = function () {
        const title = $('#notice_update_title').val();
        const content = $('#notice_update_content').val();
        const topyn = $('#notice_update_topyn').val();

        if (!title || !content) {
            alert('제목과 내용을 모두 입력해주세요!!!');
            return;
        }

        $.ajax({
            url: '/admin/notice/update',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                NOTICE_ID: selectedData.NOTICE_ID,
                TITLE: title,
                CONTENT: content,
                TOP_YN: topyn === 'true'
            }),
            success: function (response) {
                alert('공지사항이 성공적으로 수정되었습니다!!!');
                $('#notice-update').modal('hide');
                location.reload();
            },
            error: function (error) {
                alert('공지사항 수정에 실패했습니다!!!');
                console.error('에러:', error);
            }
        });
    };

    // 공지사항 삭제
    $('#DELETE_NOTICE').on('click', function (e) {
        if (!selectedData) {
            e.preventDefault();
            alert('삭제할 공지사항을 선택해주세요!!!');
            return;
        }

        if (confirm('선택한 공지사항을 삭제하시겠습니까???')) {
            $.ajax({
                url: '/admin/notice/delete',
                method: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({ NOTICE_ID: selectedData.NOTICE_ID }),
                success: function (response) {
                    alert('공지사항이 성공적으로 삭제되었습니다!!!');
                    location.reload();
                },
                error: function (error) {
                    alert('공지사항 삭제에 실패했습니다!!!');
                    console.error('에러:', error);
                }
            });
        }
    });

    // 공지사항 수정 모달이 닫힐 때 초기화
    $('#notice-update').on('hidden.bs.modal', function () {
        $(this).find('input, textarea, select').val('');
        selectedData = null;
        location.reload();
    });

    // 도서 관리 - 수정 버튼 클릭시 모달에 바인딩
    $('#update_book').on('click', function (e) {
        if (!selectedData) {
            e.preventDefault();
            alert('수정할 도서를 선택해주세요!!!');
            return;
        }

        // 기본 정보 바인딩
        $('#u_book_title').val(selectedData.BOOK_TITLE);
        $('#u_author').val(selectedData.AUTHOR);
        $('#u_isbn').val(selectedData.ISBN);
        $('#u_isbn_add').val(selectedData.ISBN_ADD);
        $('#u_publisher').val(selectedData.PUBLISHER);
        $('#u_price').val(selectedData.PRICE);
        $('#u_page').val(selectedData.PAGE);
        $('#u_ebook_yn').val(selectedData.EBOOK_YN);
        $('#u_grade').val(selectedData.GRADE);
        $('#u_publish_date').val(selectedData.PUBLISH_DATE);
        $('#u_cover_url').val(selectedData.COVER_URL);
        $('#u_contents_list').val(selectedData.CONTENTS_LIST);
        $('#u_description').val(selectedData.DESCRIPTION);
        $('#u_summary').val(selectedData.SUMMARY);
        $('#u_recommendation').val(selectedData.RECOMMENDATION);

        // 태그 선택 설정
        const tags = selectedData.TAGS.split(',');
        const choicesInstance = new Choices('#u_book_tag');
        choicesInstance.setValue(tags);

        // 개발자 추천 정보 바인딩 - 이제 별도 요청 없이 selectedData에서 바로 사용
        if (selectedData.DEV_RECOMMENDS && selectedData.DEV_RECOMMENDS.length > 0) {
            selectedData.DEV_RECOMMENDS.forEach((rec, index) => {
                if (index > 0) {
                    $('#u_add-recommender').click();  // 필요한 만큼 추천인 폼 추가
                }
                $(`#u_devID${index + 1}`).val(rec.DEV_ID);
                $(`#u_name${index + 1}`).val(rec.DEV_NAME);
                $(`#u_reason${index + 1}`).val(rec.DEV_RECOMMEND_REASON);
            });
        }

        // 수정 모달 표시
        const updateModal = new bootstrap.Modal(document.getElementById('book-update'));
        updateModal.show();
    });

    // 도서 수정 실행
    $('#update_book_btn').on('click', function () {
        // 폼 데이터 수집
        const bookData = {
            BOOK_ID: selectedData.BOOK_ID,
            BOOK_TITLE: $('#u_book_title').val(),
            AUTHOR: $('#u_author').val(),
            ISBN: $('#u_isbn').val(),
            ISBN_ADD: $('#u_isbn_add').val(),
            PUBLISHER: $('#u_publisher').val(),
            PRICE: parseInt($('#u_price').val()),
            PAGE: parseInt($('#u_page').val()),
            EBOOK_YN: $('#u_ebook_yn').val(),
            GRADE: $('#u_grade').val(),
            PUBLISH_DATE: $('#u_publish_date').val(),
            COVER_URL: $('#u_cover_url').val(),
            CONTENTS_LIST: $('#u_contents_list').val(),
            DESCRIPTION: $('#u_description').val(),
            SUMMARY: $('#u_summary').val(),
            RECOMMENDATION: $('#u_recommendation').val(),
            TAGS: Array.from($('#u_book_tag option:selected')).map(opt => opt.value).join(',')
        };

        // 개발자 추천 정보 수집
        const devRecommends = [];
        $('#u_recommenderContainer > div').each(function (index) {
            const groupIndex = Math.floor(index / 3) + 1;
            if ($(`#u_devID${groupIndex}`).val()) {
                devRecommends.push({
                    DEV_ID: $(`#u_devID${groupIndex}`).val(),
                    RECOMMEND_REASON: $(`#u_reason${groupIndex}`).val()
                });
            }
        });

        // 서버로 전송
        $.ajax({
            url: '/admin/book/update',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                book: bookData,
                devRecommends: devRecommends
            }),
            success: function (response) {
                alert('도서가 성공적으로 수정되었습니다!!!');
                $('#book-update').modal('hide');
                location.reload();
            },
            error: function (error) {
                alert('도서 수정에 실패했습니다!!!');
                console.error('에러:', error);
            }
        });
    });

    // 수정 모달이 닫힐 때 초기화
    $('#book-update').on('hidden.bs.modal', function () {
        $(this).find('input, textarea, select').val('');
        selectedData = null;
        location.reload();
    });
});

// 개발자 관리 - 개발자 삭제
$('#DELETE_DEV').on('click', function (e) {
    if (!selectedData) {
        e.preventDefault();
        alert('삭제할 개발자를 선택해주세요.');
        return;
    }
    if (confirm('선택한 개발자를 삭제하시겠습니까?')) {
        $.ajax({
            url: '/admin/dev/delete',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ ID: selectedData.ID }),
            success: function (response) {
                alert("개발자 삭제에 성공했습니다");
                $('#dev-add-update').modal('hide');
                location.reload();
            },
            error: function (error) {
                alert("개발자 삭제 실패");
            }
        });
    }
});

// 도서 삭제
// 공지사항 삭제
$('#DELETE_BOOK').on('click', function (e) {
    if (!selectedData) {
        e.preventDefault();
        alert('삭제할 도서를 선택해주세요!!!');
        return;
    }

    if (confirm('선택한 도서를 삭제하시겠습니까???')) {
        $.ajax({
            url: '/admin/book/delete',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ BOOK_ID: selectedData.BOOK_ID }),
            success: function (response) {
                alert('도서가 성공적으로 삭제되었습니다!!!');
                location.reload();
            },
            error: function (error) {
                alert('도서 삭제에 실패했습니다!!!');
                console.error('에러:', error);
            }
        });
    }
});