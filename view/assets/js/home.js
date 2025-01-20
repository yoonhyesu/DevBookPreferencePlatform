document.addEventListener('DOMContentLoaded', function () {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl)
    })
});

// 이미지 URL 처리 함수 수정
function getImageUrl(path) {
    if (!path) return '/assets/images/no-image.jpg';

    // ProfileImagePath: 접두어 제거
    path = path.replace('ProfileImagePath:', '');

    // URL 디코딩 처리
    try {
        path = decodeURIComponent(path);
    } catch (e) {
        console.error("URL 디코딩 실패:", e);
    }

    // 경로가 /uploads로 시작하지 않으면 추가
    if (!path.startsWith('/uploads')) {
        path = '/uploads' + path;
    }

    // 백슬래시를 슬래시로 변환
    path = path.split('\\').join('/');

    // 중복된 슬래시 제거
    path = path.replace(/\/+/g, '/');

    console.log("최종 이미지 경로:", path);
    return path;
}

// 개발자 카드 생성 함수에 디버깅 로그 추가
function createDevCard(dev) {
    console.log("개발자 데이터:", dev);
    const imagePath = getImageUrl(dev.ProfileImagePath);
    console.log("변환된 이미지 경로:", imagePath);

    return `
        <div class="col">
            <div class="card shadow-sm h-100">
                <img src="${imagePath}" class="card-img-top dev-profile-img" 
                     alt="${dev.DevName}" style="height: 225px; object-fit: cover;"
                     onerror="this.src='/assets/images/no-image.jpg'">
                <div class="card-body">
                    <h5 class="card-title">${dev.DevName}</h5>
                    <p class="card-text">${dev.DevDetailName || ''}</p>
                    <div class="d-flex justify-content-between align-items-center">
                        <small class="text-body-secondary">경력: ${dev.DevHistory || 'N/A'}</small>
                    </div>
                </div>
            </div>
        </div>
    `;
}