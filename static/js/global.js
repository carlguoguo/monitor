$(function() {
    //$(document).foundation();

    // Hide any messages after a few seconds
    hideFlash();
    $("#api-model").on("hidden.bs.modal", function() {
        location.reload();
    });
    $('#api-model').on('show.bs.modal', function(event) {
        var button = $(event.relatedTarget); // Button that triggered the modal
        var modal = $(this);
        modal.find('.modal-title').append(button.data('title'));
        modal.find('.modal-body .alias').val(button.data('alias'));
        var interval_time = button.data('interval');
        modal.find('.modal-body .url').val(button.data('url'));
        modal.find('.modal-body .receivers').val(button.data('alert'));
        modal.find('.modal-body .timeout').val(button.data('timeout'));
        modal.find('.modal-body .fail_threshold').val(button.data('failmax'));
        var action = button.data('action');
        if (action == "update"){
            modal.find('form').attr('action', '/api/update/' + button.data('id'));
        } else {
            modal.find('form').attr('action', '/api/create');
        }
        var handle = $("#custom-handle");
        $("#slider").slider({
            range: "max",
            min: 1,
            max: 10,
            value: interval_time,
            create: function() {
                handle.text($(this).slider("value") + " 分钟");
            },
            slide: function(event, ui) {
                $("#interval_time").val(ui.value);
                handle.text(ui.value + " 分钟");
            }
        });
        $("#interval_time").val($("#slider").slider("value"));
    });
    $(".clickable-row").click(function() {
        window.location = $(this).data("href");
    });
});

function hideFlash(rnum) {
    if (!rnum) rnum = '0';

    _.delay(function() {
        $('.alert-box-fixed' + rnum).fadeOut(300, function() {
            $(this).css({
                "visibility": "hidden",
                display: 'block'
            }).slideUp();

            var that = this;

            _.delay(function() {
                that.remove();
            }, 400);
        });
    }, 4000);
}

function showFlash(obj) {
    $('#flash-container').html();
    $(obj).each(function(i, v) {
        var rnum = _.random(0, 100000);
        var message = '<div id="flash-message" class="alert-box-fixed' +
            rnum + ' alert-box-fixed alert alert-dismissible ' + v.cssclass + '">' +
            '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>' +
            v.message + '</div>';
        $('#flash-container').prepend(message);
        hideFlash(rnum);
    });
}

function flashError(message) {
    var flash = [{
        Class: "alert-danger",
        Message: message
    }];
    showFlash(flash);
}

function flashSuccess(message) {
    var flash = [{
        Class: "alert-success",
        Message: message
    }];
    showFlash(flash);
}

function flashNotice(message) {
    var flash = [{
        Class: "alert-info",
        Message: message
    }];
    showFlash(flash);
}

function flashWarning(message) {
    var flash = [{
        Class: "alert-warning",
        Message: message
    }];
    showFlash(flash);
}
