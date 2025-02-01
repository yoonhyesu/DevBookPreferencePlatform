$(document).ready(function () {
    // $('#tagInput').on("keyup", function () {
    //     var value = $(this).val().toLowerCase();
    //     $('#tagList li button').filter(function () {
    //         $(this).toggle($(this).text().toLowerCase().indexOf(value) > -1);
    //     }).css('border', '2px solid red');
    // })


    $('#tagInput').on("keyup", function () {
        var value = $(this).val().toLowerCase();
        $('#tagList li button').each(function () {
            if (value == "") {
                $(this).css({
                    'border': '3px solid white',
                    'background': 'white',
                    'transition': 'all 0.3s ease'
                });
            } else {
                var matches = $(this).text().toLowerCase().indexOf(value) > -1;
                if (matches) {
                    $(this).css({
                        'border': '3px solid yellow',
                        'background': 'linear-gradient(45deg, #fff, #ffffd0)',
                        'transition': 'all 0.3s ease'
                    });
                } else {
                    $(this).css({
                        'border': '3px solid white',
                        'background': 'white',
                        'transition': 'all 0.3s ease'
                    });
                }
            }
        });
    });
});