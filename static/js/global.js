$(function() {
    //$(document).foundation();

	// Hide any messages after a few seconds
    hideFlash();
	$('#updateAPI').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget) // Button that triggered the modal
		var title = button.data('whatever') // Extract info from data-* attributes
		// If necessary, you could initiate an AJAX request here (and then do the updating in a callback).
		// Update the modal's content. We'll use jQuery here, but you could use a data binding library or other methods instead.
		var modal = $(this)
		modal.find('.modal-body .alias').val(button.data('alias'))
		modal.find('.modal-body .interval_time').val(button.data('interval'))
		modal.find('.modal-body .url').val(button.data('url'))
		modal.find('.modal-body .receivers').val(button.data('alert'))
		modal.find('.modal-body .timeout').val(button.data('timeout'))
		modal.find('.modal-body .fail_threshold').val(button.data('failmax'))
		modal.find('form').attr('action', '/api/update/'+button.data('id'));
	})
});

function hideFlash(rnum)
{
    if (!rnum) rnum = '0';

    _.delay(function() {
        $('.alert-box-fixed' + rnum).fadeOut(300, function() {
            $(this).css({"visibility":"hidden",display:'block'}).slideUp();

            var that = this;

            _.delay(function() { that.remove(); }, 400);
        });
    }, 4000);
}

function showFlash(obj)
{
    $('#flash-container').html();
    $(obj).each(function(i, v) {
        var rnum = _.random(0, 100000);
		var message = '<div id="flash-message" class="alert-box-fixed'
		+ rnum + ' alert-box-fixed alert alert-dismissible '+v.cssclass+'">'
		+ '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'
		+ v.message + '</div>';
        $('#flash-container').prepend(message);
        hideFlash(rnum);
    });
}

function flashError(message) {
	var flash = [{Class: "alert-danger", Message: message}];
	showFlash(flash);
}

function flashSuccess(message) {
	var flash = [{Class: "alert-success", Message: message}];
	showFlash(flash);
}

function flashNotice(message) {
	var flash = [{Class: "alert-info", Message: message}];
	showFlash(flash);
}

function flashWarning(message) {
	var flash = [{Class: "alert-warning", Message: message}];
	showFlash(flash);
}
