$(document).ready(function () {
    const grade = $('#grade').val();
    $.ajax({
        url: "/api/recommend/step/" + grade,
        method: "GET",
        success: function (data) {
            console.log(data);
        }
    })
})