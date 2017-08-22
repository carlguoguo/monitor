$(function() {
    //$(document).foundation();

    // Hide any messages after a few seconds
    hideFlash();
    $("#api-model").on("hidden.bs.modal", function() {
        location.reload();
    });
    $('#api-model').on('show.bs.modal', function(event) {
        var button = $(event.relatedTarget); // Button that triggered the modal
        var title = button.data('whatever'); // Extract info from data-* attributes
        // If necessary, you could initiate an AJAX request here (and then do the updating in a callback).
        // Update the modal's content. We'll use jQuery here, but you could use a data binding library or other methods instead.
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
                handle.text($(this).slider("value"));
            },
            slide: function(event, ui) {
                $("#interval_time").val(ui.value);
                handle.text(ui.value);
            }
        });
        $("#interval_time").val($("#slider").slider("value"));
    })
    // $('#updateAPI').on('show.bs.modal', function (event) {
    //     var button = $(event.relatedTarget); // Button that triggered the modal
    // 	var title = button.data('whatever'); // Extract info from data-* attributes
    // 	// If necessary, you could initiate an AJAX request here (and then do the updating in a callback).
    // 	// Update the modal's content. We'll use jQuery here, but you could use a data binding library or other methods instead.
    // 	var modal = $(this);
    // 	modal.find('.modal-body .alias').val(button.data('alias'));
    //     var interval_time = button.data('interval');
    // 	modal.find('.modal-body .url').val(button.data('url'));
    // 	modal.find('.modal-body .receivers').val(button.data('alert'));
    // 	modal.find('.modal-body .timeout').val(button.data('timeout'));
    // 	modal.find('.modal-body .fail_threshold').val(button.data('failmax'));
    // 	modal.find('form').attr('action', '/api/update/'+button.data('id'));
    //     var handle = $( "#updated-custom-handle" );
    //     $("#updated_slider").slider({
    //         range: "max",
    //         min: 1,
    //         max: 10,
    //         value: interval_time,
    //         create: function() {
    //             handle.text($(this).slider("value"));
    //         },
    //         slide: function(event, ui) {
    //             $("#updated_interval_time").val(ui.value);
    //             handle.text(ui.value);
    //         }
    //     });
    //     $("#updated_interval_time").val($("#updated_slider").slider("value"));
    // })
    //
    // $('#createAPI').on('show.bs.modal', function (event) {
    //     var button = $(event.relatedTarget); // Button that triggered the modal
    // 	var title = button.data('whatever'); // Extract info from data-* attributes
    // 	// If necessary, you could initiate an AJAX request here (and then do the updating in a callback).
    // 	// Update the modal's content. We'll use jQuery here, but you could use a data binding library or other methods instead.
    // 	var modal = $(this);
    // 	modal.find('.modal-body .alias').val(button.data('alias'));
    //     var interval_time = button.data('interval');
    // 	modal.find('.modal-body .url').val(button.data('url'));
    // 	modal.find('.modal-body .receivers').val(button.data('alert'));
    // 	modal.find('.modal-body .timeout').val(button.data('timeout'));
    // 	modal.find('.modal-body .fail_threshold').val(button.data('failmax'));
    // 	modal.find('form').attr('action', '/api/update/'+button.data('id'));
    //     var handle = $( "#created-custom-handle" );
    //     $("#created_slider").slider({
    //         range: "max",
    //         min: 1,
    //         max: 10,
    //         value: interval_time,
    //         create: function() {
    //             handle.text($(this).slider("value"));
    //         },
    //         slide: function(event, ui) {
    //             $("#created_interval_time").val(ui.value);
    //             handle.text(ui.value);
    //         }
    //     });
    //     $("#created_interval_time").val($("#created_slider").slider("value"));
    // })
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
