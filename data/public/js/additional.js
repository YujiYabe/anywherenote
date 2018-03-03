"use strict";


// ページロード後処理
$(function () {
    makePageList();


});




function makePageList() {

    $('#parent_note_table').empty();
    $('#page_star_div').empty();
    $('#table_parent').empty();

    var json_return_value = JSON.parse($('#source_dataset_list').text());

    // statuscode == '1' 保存先が一つもない場合、ノート追加を表示、ページ追加を非表示
    if (json_return_value["RtnCode"] == '1') {
        // if (true) {
        $('#flexbox_page_pane').hide();
        $('#flexbox_note_pane').show();
        $('#flexbox_left_pane').hide();
        // $('#left_pane_head').hide();
        return false;
    } else {
        $('#flexbox_note_pane').hide();

    }

    var select_note_id = json_return_value["SlctPst"].NoteID
    var select_page_id = json_return_value["SlctPst"].PageID

    // console.log(select_note_id + ":" + select_page_id);
    // console.log("=================");

    var parent_note_table = $('#parent_note_table');

    // var note_table = $('<div>'); parent_note_table.append(note_table); note_table.addClass('container'); //table.attr('id', note_address); table.attr('data-address', note_address); 
    var note_tr = $('<div>'); parent_note_table.append(note_tr); note_tr.addClass('note_tr');




    $.each(json_return_value["DataSet"], function (index, val) {
        var note_id = json_return_value["DataSet"][index]["NoteID"];

        var note_star = json_return_value["DataSet"][index]["NoteStar"];

        var note_name = json_return_value["DataSet"][index]["NoteName"];
        var note_address = json_return_value["DataSet"][index]["NoteAddress"];
        // var temp_note_updatetime = json_return_value["DataSet"][index]["NoteDBUpdateTime"];
        var note_updatetime = moment(json_return_value["DataSet"][index]["NoteUpdateTime"]).format('YYYY/MM/DD HH:mm:ss');

        var page_list = json_return_value["DataSet"][index]["list"];

        // $('#createPagebutton').attr('name', note_address);

        var note_tr = $('<div>'); parent_note_table.append(note_tr); note_tr.addClass('note_tr');

        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('note_td_1'); note_td.text("変更"); note_td.addClass("btn btn-info"); note_td.attr('onclick', 'updateNote(this);');
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('note_td_2');


        // 右ペイン各ノートのスター文字列
        var right_pane_note_star = convertStar(note_star, true); note_td.append(right_pane_note_star);

        // 左ペイン各ノートのスター文字列
        var left_pane_note_star = convertStar(note_star, false);


        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('note_td_3');
        var note_input = $('<input>'); note_td.append(note_input); note_input.addClass("form-control  input-sm");
        note_input.val(note_name);
        note_input.attr('data-note_id', note_id);
        note_input.attr('data-star_int', note_star);

        var message_span = $('<span>'); note_td.append(message_span); message_span.hide(); message_span.text('更新しました');

        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('note_td_4'); note_td.text(note_address);
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('note_td_5'); note_td.text("削除"); note_td.addClass("btn btn-danger"); note_td.attr('onclick', 'deleteNote(this);');

        var element = $('#table_parent');

        var hr = $('<hr>'); element.append(hr);

        //---------------

        var span = $('<div>'); element.append(span); span.addClass("btn btn-primary btn-lg notedb_button"); span.attr('data-target_table', note_id); span.attr('onclick', 'switchShowHideDataList(this);');
        var child_span = $('<div>'); span.append(child_span); child_span.addClass(""); child_span.text(note_name);
        var child_span = $('<div>'); span.append(child_span); child_span.addClass("notedb_datetime"); child_span.append(left_pane_note_star);
        var child_span = $('<div>'); span.append(child_span); child_span.addClass("notedb_address"); child_span.text(note_address);


        //---------------
        var parent_div = $('<div>'); element.append(parent_div); parent_div.attr('id', note_id); parent_div.hide();

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
        span.attr('onclick', "createPage(this);");

        var h5 = $('<h5>'); parent_div.append(h5);

        var table = $('<table>'); parent_div.append(table);
        table.attr('id', note_address);
        table.attr('data-note_id', note_id);
        table.attr('data-note_name', note_name);
        table.attr('data-note_address', note_address);

        var thead = $('<thead>'); table.append(thead);

        var tr = $('<tr>'); thead.append(tr);
        var th = $('<th>'); tr.append(th); th.hide(); th.text('summary');
        var th = $('<th>'); tr.append(th); th.hide(); th.text('body');

        var tbody = $('<tbody>'); table.append(tbody);


        for (var item in page_list) {
            var page_star = page_list[item]['PageStar'];

            var left_pane_page_star = convertStar(page_star);


            // var temp_date = Date.parse(page_list[item]['UpdatedAt']);
            // var updateDateTime = moment(temp_date).format('YYYY/MM/DD HH:mm:ss');

            var tr = $('<tr>'); tbody.append(tr); tr.attr('data-note_id', note_id); tr.addClass('page_item'); tr.attr('data-page_id', page_list[item]['PageID']);

            var td = $('<td>'); tr.append(td);
            var parent_div = $('<div>'); td.append(parent_div); parent_div.addClass('btn btn-info note_object'); parent_div.attr('onclick', 'showDataToRightPane(this)');

            var temp_page_body = page_list[item]['PageBody'].replace(new RegExp('\/\/note_id\/\/', "g"), note_id);

            parent_div.attr('data-update_time', page_list[item]['PageUpdate']);
            parent_div.attr('data-page_id', page_list[item]['PageID']);
            parent_div.attr('data-page_title', page_list[item]['PageTitle']);
            parent_div.attr('data-page_body', temp_page_body);
            parent_div.attr('data-page_star', page_list[item]['PageStar']);


            var div = $('<div>'); parent_div.append(div); div.addClass('left_pane_page_object'); div.append(left_pane_page_star);
            var div = $('<div>'); parent_div.append(div); div.addClass('left_pane_page_object'); div.text(page_list[item]['PageTitle']);

            var td = $('<td>'); tr.append(td); td.hide(); td.html(temp_page_body);


            // 選択済のページを表示
            if (
                select_note_id == note_id &&
                select_page_id == page_list[item]['PageID']
            ) {
                // tr.addClass("currentItem");
                $('#note_id').text(note_id);
                $('#note_name').text(note_name);
                $('#note_address').text(note_address);

                $('#post_note_address').val(note_address);
                $('#post_note_id').val(note_id);
                $('#post_page_id').val(page_list[item]['PageID']);

                $('#page_id').text(page_list[item]['PageID']);
                $('#page_title').val(page_list[item]['PageTitle']);

                var right_pane_page_star = convertStar(page_list[item]['PageStar'], true);

                $('#page_star_div').append(right_pane_page_star);

                var temp_page_body = page_list[item]['PageBody'].replace(new RegExp('\/\/note_id\/\/', "g"), note_id);

                $('#page_body').html(temp_page_body);
                $('#update_time').text(page_list[item]['PageUpdate']);

                $('#flexbox_page_pane').show();

                parent_div.addClass("currentItem");

            }

        }

    });


    $('input#search').quicksearch('table tbody tr');
    // $(".fileDownload").click(function () {
    $('.fileDownload').on('click', function () {

        var fileurl = $(this).attr('href');
        if (!confirm($(this).text() + 'をダウンロードしますか？')) {
            /* キャンセルの時の処理 */
            return false;
        }
        var childWindow = window.open('about:blank');
        $.ajax({
            type: 'GET',
            url: fileurl,
        }).done(function (jqXHR) {
            childWindow.location.href = fileurl;

            var id = setInterval(function () {
                if (true) {
                    childWindow.close();
                    clearInterval(id);
                }
            }, 10);

        }).fail(function (jqXHR) {
            childWindow.close();
        });
    });
} // =======================================




// スターをクリックした際のレートの変更
function changeRateStar(obj) {

    // $(obj).parent().prevAll().eq(0).text($(obj).attr('data-number'));
    $(obj).parent().parent().nextAll().eq(0).children('input').attr('data-star_int', $(obj).attr('data-number'));

} // =======================================



// ノートをクリックした際にページの表示・非表示
function switchShowHideDataList(obj) {
    var target_table = $(obj).attr('data-target_table');
    $('#' + target_table).toggle();

} // =======================================


// ----------------------------------------------------
// 右ペインの内容変更 テキストエリア⇔アップロードエリア
function switch_rght_body() {

    if ($('#upload_zone').css('display') == 'none') {
        // テキストエリア⇒アップロードエリアの場合、現在のテキストを保存してからアップロード処理
        updatePage();
        $('#page_body').css('display', 'none')
        $('#upload_zone').css('display', '')

    } else {
        $('#page_body').css('display', '')
        $('#upload_zone').css('display', 'none')
    }

} // ----------------------------------------------------


// ----------------------------------------------------
// 右ペインの変更 ノート編集⇔ページ編集
function switchRightPane(pane_name) {
    $('.right_pane_content').hide();
    $('#' + pane_name).show();
} // ----------------------------------------------------

// ----------------------------------------------------
// 星の数字を文字列に変換
function convertStar(star_number, isEnableChange) {

    var parent_span = $('<span>');

    for (var starint = 1; starint <= 3; starint++) {
        var span = $('<span>'); parent_span.append(span);

        if (isEnableChange) {
            span.attr('data-number', starint); span.addClass('ratingstar'); span.attr('onclick', 'changeRateStar(this)');
            span.attr('onmouseover', 'hoverStar(this)');
            span.attr('onmouseout', 'blurStar(this)');
        }

        if (starint <= star_number) {
            span.text('★');
        } else {
            span.text('☆');
        }
    }

    return parent_span;


} // ----------------------------------------------------


// ----------------------------------------------------
// クリックしたページを右ペインに表示
function showDataToRightPane(obj) {
    $(".currentItem").removeClass("currentItem");

    switchRightPane('flexbox_page_pane');
    $('#page_star_div').empty();


    //            div    td       tr       tbody    table
    var note_info = $(obj).parent().parent().parent().parent();

    var note_id = note_info.attr('data-note_id');
    var note_name = note_info.attr('data-note_name');
    var note_address = note_info.attr('data-note_address');
    // var note_id = $(obj).parent().parent().parent().parent().attr('data-note_id');
    // var note_name = $(obj).parent().parent().parent().parent().attr('data-note_name');
    // var note_address = $(obj).parent().parent().parent().parent().attr('data-note_address');


    var update_time = $(obj).attr('data-update_time');
    var page_id = $(obj).attr('data-page_id');
    var page_title = $(obj).attr('data-page_title');
    var page_body = $(obj).attr('data-page_body');
    var page_star = $(obj).attr('data-page_star');

    // var update_time = $(obj).parent().nextAll().eq(0).text();
    // var page_id = $(obj).parent().nextAll().eq(1).text();
    // var page_title = $(obj).parent().nextAll().eq(2).text();
    // var page_body = $(obj).parent().nextAll().eq(3).text();
    // var page_star = $(obj).parent().nextAll().eq(4).text();

    // console.log(page_star);

    // console.log(update_time);
    $('#update_time').text(update_time);

    $('#note_id').text(note_id);
    $('#note_name').text(note_name);
    $('#note_address').text(note_address);

    $('#page_id').text(page_id);
    $('#page_title').val(page_title);
    $('#page_body').html(page_body);


    $('#post_note_address').val(note_address);
    $('#post_note_id').val(note_id);
    $('#post_page_id').val(page_id);

    var right_pane_page_star = convertStar(page_star, true);

    $('#page_star_div').append(right_pane_page_star);

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
    //     console.log('IDの重複はありません。');
    // }
    // //ID重複チェック ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝

    // console.log(document.querySelectorAll('[id]'));


} // =======================================



function hoverStar(Obj) {
    $(Obj).siblings().text('☆');
    $(Obj).prevAll().text('★');
    $(Obj).text('★');
}

function blurStar(Obj) {
    var selected_int = $(Obj).parent().parent().nextAll().eq(0).children('input').attr('data-star_int');

    var parent = $(Obj).parent();

    for (var starint = 1; starint <= 3; starint++) {
        var targetspan = parent.find('[data-number=' + starint + ']');

        if (starint <= selected_int) {
            targetspan.text('★');
        } else {
            targetspan.text('☆');
        }
    }

}
