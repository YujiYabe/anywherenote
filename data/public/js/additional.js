"use strict";

$(function () {

    $("#line_separate").hover(
        function (e) {
            $(this).css({ background: "blue" });
            $("#line_separate").css("width", "15px")
        },
        function (e) {
            $("#line_separate").css("width", "3px")
        }
    );

    $("#line_separate").draggable({
        axis: "x",
        start: function (event, ui) {
            $(this).css({ background: "blue" });
        },
        drag: function (event, ui) {
            // console.log(event , ui);
            // console.log(' top: ' + ui.position.top + ' left: ' + ui.position.left);
            $(this).css({ background: "blue" });
            $(".left_pane").css("width", ui.position.left)
            // $( "BODY" ).css("padding-left", ui.position.left + 5)
            $(".right_pane").css("padding-left", ui.position.left + 5)
        },

        stop: function (event, ui) {
            $(this).css({ background: "gray" });
        }
    });
});


$(function () {
    makePageList();

    $('input#search').quicksearch('table tbody tr');

});

function switch_rght_body() {

    if ($('#upload_zone').css('display') == 'none') {
        $('#page_body').css('display', 'none')
        $('#upload_zone').css('display', '')
    } else {
        $('#page_body').css('display', '')
        $('#upload_zone').css('display', 'none')
    }

}


function switchRightPane(pane_name) {

    $('.right_pane_content').hide();
    $('#' + pane_name).show();
}

function makePageList() {

    $('#parent_note_table').empty();
    $('#table_parent').empty();

    var json_return_value = JSON.parse($('#source_return_value').text());

    // statuscode == '1' 保存先が一つもない場合、ノート追加を表示、ページ追加を非表示
    if (json_return_value["RtnCode"] == '1') {
        $('#edit_page').hide();
        $('#add_note').show();
        $('#left_pane_head').hide();
        return false;
    } else {
        $('#add_note').hide();

    }

    var select_note_id = json_return_value["SlctPst"].NoteID
    var select_page_id = json_return_value["SlctPst"].PageID

    // console.log(select_note_id + ":" + select_page_id);
    // console.log("=================");

    var parent_note_table = $('#parent_note_table');

    var note_table = $('<div>'); parent_note_table.append(note_table); note_table.addClass('container'); //table.attr('id', note_address); table.attr('data-address', note_address); 
    var note_tr = $('<div>'); note_table.append(note_tr); note_tr.addClass('row');
    var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-1'); note_td.text("追加");
    var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-4'); note_td.text("ノート名");
    var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-4'); note_td.text("ノート格納パス");
    var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-1'); note_td.text("削除");


    var note_tbody = $('<tbody>'); note_table.append(note_tbody);

    $.each(json_return_value["DataSet"], function (index, val) {
        var note_id = json_return_value["DataSet"][index]["NoteDBID"];
        var note_name = json_return_value["DataSet"][index]["NoteDBName"];
        var note_address = json_return_value["DataSet"][index]["NoteDBAddress"];
        // var temp_note_updatetime = json_return_value["DataSet"][index]["NoteDBUpdateTime"];
        var note_updatetime = moment(json_return_value["DataSet"][index]["NoteDBUpdateTime"]).format('YYYY/MM/DD HH:mm:ss');

        var page_list = json_return_value["DataSet"][index]["list"];

        $('#addPagebutton').attr('name', note_address);


        var h5 = $('<h5>'); note_table.append(h5);
        var note_tr = $('<div>'); note_table.append(note_tr); note_tr.addClass('row');

        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-1'); note_td.text("変更"); note_td.addClass("btn btn-info"); note_td.attr('onclick', 'updateNote(this);');
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-4');
        var note_input = $('<input>'); note_td.append(note_input); note_input.val(note_name); note_input.attr('data-note_id', note_id); note_input.addClass("form-control");
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-4'); note_td.text(note_address);
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-1'); note_td.text("削除"); note_td.addClass("btn btn-danger"); note_td.attr('onclick', 'deleteNote(this);');

        var element = $('#table_parent');


        //---------------

        var span = $('<div>'); element.append(span); span.addClass("btn btn-primary btn-lg notedb_button"); span.attr('data-target_table', note_id); span.attr('onclick', 'switchShowHideDataList(this);');
        var child_span = $('<div>'); span.append(child_span); child_span.addClass("notedb_datetime"); child_span.text(note_updatetime);
        var child_span = $('<div>'); span.append(child_span); child_span.addClass(""); child_span.text(note_name);
        var child_span = $('<div>'); span.append(child_span); child_span.addClass("notedb_address"); child_span.text(note_address);


        //---------------
        var parent_div = $('<div>'); element.append(parent_div); parent_div.attr('id', note_id); parent_div.hide();
        // if (index != 0) {
        // console.log(index);
        // console.log(select_note_id);

        if (note_id == select_note_id) {
            parent_div.show();
        }

        var h5 = $('<h5>'); parent_div.append(h5);
        //---------------

        var span = $('<span>'); parent_div.append(span);
        span.addClass("btn  btn-secondary ");
        span.text('ページ追加');
        span.attr('data-select_note_id', note_id);
        span.attr('data-address', note_address);
        span.attr('onclick', "addPage(this);");

        var h5 = $('<h5>'); parent_div.append(h5);

        var table = $('<table>'); parent_div.append(table);
        table.attr('id', note_address);
        table.attr('data-note_id', note_id);
        table.attr('data-note_name', note_name);
        table.attr('data-note_address', note_address);
        table.addClass('table table-bordered table-hover');

        var thead = $('<thead>'); table.append(thead);

        var tr = $('<tr>'); thead.append(tr);
        var th = $('<th>'); tr.append(th); th.hide(); th.text('summary');
        var th = $('<th>'); tr.append(th); th.hide(); th.text('updatetime');
        var th = $('<th>'); tr.append(th); th.hide(); th.text('ID');
        var th = $('<th>'); tr.append(th); th.hide(); th.text('title');
        var th = $('<th>'); tr.append(th); th.hide(); th.text('body');

        var tbody = $('<tbody>'); table.append(tbody);


        for (var item in page_list) {

            var temp_date = Date.parse(page_list[item]['UpdatedAt']);
            var updateDateTime = moment(temp_date).format('YYYY/MM/DD HH:mm:ss');

            var tr = $('<tr>'); tbody.append(tr); tr.attr('data-note_id', note_id); tr.attr('data-page_id', page_list[item]['ID']);

            var td = $('<td>'); tr.append(td); td.attr('onclick', 'showDataToRightPane(this)');
            var div = $('<div>'); td.append(div); div.text(updateDateTime);
            var div = $('<div>'); td.append(div); div.text(page_list[item]['page_title']);

            var td = $('<td>'); tr.append(td); td.hide(); td.text(updateDateTime);
            var td = $('<td>'); tr.append(td); td.hide(); td.text(page_list[item]['ID']);
            var td = $('<td>'); tr.append(td); td.hide(); td.text(page_list[item]['page_title']);

            var temp_page_body = page_list[item]['page_body'].replace(new RegExp('\/\/note_id\/\/', "g"), note_id);

            var td = $('<td>'); tr.append(td); td.hide(); td.text(temp_page_body);


            // 選択済のページを表示
            if (
                select_note_id == note_id &&
                select_page_id == page_list[item]['ID']
            ) {
                tr.addClass("currentItem");
                $('#note_id').text(note_id);
                $('#note_name').text(note_name);
                $('#note_address').text(note_address);

                $('#post_note_address').val(note_address);
                $('#post_note_id').val(note_id);
                $('#post_page_id').val(page_list[item]['ID']);

                $('#page_id').text(page_list[item]['ID']);
                $('#page_title').val(page_list[item]['page_title']);

                var temp_page_body = page_list[item]['page_body'].replace(new RegExp('\/\/note_id\/\/', "g"), note_id);

                $('#page_body').html(temp_page_body);
                $('.update_time').text(updateDateTime);

                $('#edit_page').show();



            }

        }
        var hr = $('<hr>'); element.append(hr);
        var hr = $('<hr>'); element.append(hr);

    });

} // =======================================

// ノートをクリックした際にページの表示・非表示
function switchShowHideDataList(obj) {
    var target_table = $(obj).attr('data-target_table');
    $('#' + target_table).toggle();


} // =======================================


// ページを右ペインに表示
function showDataToRightPane(obj) {
    $(".currentItem").removeClass("currentItem");

    switchRightPane('edit_page');


    var note_id = $(obj).parent().parent().parent().attr('data-note_id');
    var note_name = $(obj).parent().parent().parent().attr('data-note_name');
    var note_address = $(obj).parent().parent().parent().attr('data-note_address');

    // console.log(note_name);
    $('#update_time').text($(obj).nextAll().eq(0).text());
    $('#page_id').text($(obj).nextAll().eq(1).text());
    $('#page_title').val($(obj).nextAll().eq(2).text());
    $('#page_body').html($(obj).nextAll().eq(3).text());

    $('#note_address').text(note_address);
    $('#note_name').text(note_name);
    $('#note_id').text(note_id);

    $('#post_note_address').val(note_address);
    $('#post_note_id').val(note_id);
    $('#post_page_id').val($(obj).nextAll().eq(1).text());

    $(obj).addClass("currentItem");


    // //ID重複チェック ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝
    // var idArr = [];
    // var duplicateIdArr = [];
    // [].forEach.call(document.querySelectorAll('[id]'), function (elm) {
    //     var id = elm.getAttribute('id');
    //     if (idArr.indexOf(id) !== -1) {
    //         duplicateIdArr.push(id);
    //     } else {
    //         idArr.push(id);
    //     }
    // });
    // if (duplicateIdArr.length > 0) {
    //     console.log('IDの重複があります:', duplicateIdArr);
    // } else {
    //     // console.log('IDの重複はありません。');
    // }
    // //ID重複チェック ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝

    // console.log(document.querySelectorAll('[id]'));


} // =======================================



