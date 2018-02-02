$(function(){
    /* btn-sidebar */
    $('#btn-sidebar').on('click', function(){
        $('.ui.sidebar').sidebar('toggle');
    });

    /* back2top */
    $(window).on('scroll', $.throttle(250, function(){
        if($(this).scrollTop() >= 100){
            $('#back2top').fadeIn();
        } else {
            $('#back2top').fadeOut();
        }
    }));
    $('#back2top').on('click', $.throttle(250, true, function(){
        $('body,html').animate({
            scrollTop: 0
        }, 800);
    }));
});
