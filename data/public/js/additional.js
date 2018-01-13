"use strict";


$(function () {


    makePageList();

    $('input#search').quicksearch('table tbody tr');

    tinymce.init({
        selector: "#page_body",
        plugins: "autoresize",
        language: "ja",
        autoresize_bottom_margin: 1,
        font_formats: 'NotoSansMono;monospace;AkrutiKndPadmini=Akpdmi-n',
        toolbar: [ // ツールバー(2段)
            // 戻る 進む | フォーマット | 太字 斜体 | 左寄せ 中央寄せ 右寄せ 均等割付 | 箇条書き 段落番号 インデントを減らす インデント
            "undo redo | formatselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent",
            // 文字サイズ 文字色 画像 リンク
            "fontsizeselect forecolor image link"
        ],
        statusbar: false, // ステータスバーを隠す
        body_class: 'my_class'

    })




});

function switchRightPane(pane_name) {

    $('.right_pane').hide();
    $('#' + pane_name).show();
}

function makePageList() {
    $('#table_parent').empty();
    var json_return_value = JSON.parse($('#source_return_value').text());

    // statuscode == '1' 保存先が一つもない場合、ノート追加を表示、ページ追加を非表示
    if (json_return_value["key0"] == '1') {
        $('#edit_page').hide()
        $('#add_note').show()
        $('#common_item').hide()
        return
    } else {
        $('#add_note').hide()

    }

    var select_note_id = json_return_value["key2"].NoteID
    var select_page_id = json_return_value["key2"].PageID

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

    $.each(json_return_value["key1"], function (index, val) {
        var note_id = json_return_value["key1"][index]["NoteDBID"];
        var note_name = json_return_value["key1"][index]["NoteDBName"];
        var note_address = json_return_value["key1"][index]["NoteDBAddress"];
        // var temp_note_updatetime = json_return_value["key1"][index]["NoteDBUpdateTime"];
        var note_updatetime = moment(json_return_value["key1"][index]["NoteDBUpdateTime"]).format('YYYY/MM/DD HH:mm:ss');

        var page_list = json_return_value["key1"][index]["list"];

        $('#addPagebutton').attr('name', note_address);


        var h5 = $('<h5>'); note_table.append(h5);
        var note_tr = $('<div>'); note_table.append(note_tr); note_tr.addClass('row');

        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-1'); note_td.text("変更"); note_td.addClass("btn btn-info"); note_td.attr('onclick', 'updateNote(this);');
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-4'); var note_input = $('<input>'); note_td.append(note_input); note_input.val(note_name); note_input.attr('data-note_id', note_id); note_input.addClass("form-control");
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-4'); note_td.text(note_address);
        var note_td = $('<div>'); note_tr.append(note_td); note_td.addClass('col-1'); note_td.text("削除"); note_td.addClass("btn btn-danger"); note_td.attr('onclick', 'deleteNote(this);');

        var element = $('#table_parent');

        var hr = $('<hr>'); element.append(hr);
        var hr = $('<hr>'); element.append(hr);

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

        var span = $('<span>'); parent_div.append(span); span.addClass("btn  btn-secondary "); span.text('ページ追加'); span.attr('data-select_note_id', note_id); span.attr('data-address', note_address); span.attr('onclick', "addPage(this);");

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
            var td = $('<td>'); tr.append(td); td.hide(); td.text(page_list[item]['page_body']);


            // 選択済のページを表示
            if (
                select_note_id == note_id &&
                select_page_id == page_list[item]['ID']
            ) {
                tr.addClass("currentItem");
                $('#note_id').text(note_id);
                $('#note_name').text(note_name);
                $('#note_address').text(note_address);

                $('#page_id').text(page_list[item]['ID']);
                $('#page_title').val(page_list[item]['page_title']);
                $('#page_body').val(page_list[item]['page_body']);
                $('#update_time').text(updateDateTime);

                $('#edit_page').show();

            }

        }

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
    tinymce.remove('#page_body');


    var note_id = $(obj).parent().parent().parent().attr('data-note_id');
    var note_name = $(obj).parent().parent().parent().attr('data-note_name');
    var note_address = $(obj).parent().parent().parent().attr('data-address');

    console.log(note_name);
    $('#update_time').text($(obj).nextAll().eq(0).text());
    $('#page_id').text($(obj).nextAll().eq(1).text());
    $('#page_title').val($(obj).nextAll().eq(2).text());
    $('#page_body').val($(obj).nextAll().eq(3).text());

    $('#note_address').text(note_address);
    $('#note_name').text(note_name);
    $('#note_id').text(note_id);

    tinymce.init({
        selector: "#page_body",
        plugins: "autoresize",
        language: "ja",
        autoresize_bottom_margin: 1,
        // font_formats: 'Arial=arial,helvetica,sans-serif;Courier New=courier new,courier,monospace;AkrutiKndPadmini=Akpdmi-n'
        font_formats: 'NotoSansMono;monospace;AkrutiKndPadmini=Akpdmi-n',
        // font_formats: 'Consolas, Courier, Monaco, monospace'
        toolbar: [ // ツールバー(2段)
            // 戻る 進む | フォーマット | 太字 斜体 | 左寄せ 中央寄せ 右寄せ 均等割付 | 箇条書き 段落番号 インデントを減らす インデント
            "undo redo | formatselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent",
            // 文字サイズ 文字色 画像 リンク
            "fontsizeselect forecolor image link"
        ],
        statusbar: false, // ステータスバーを隠す
        body_class: 'my_class'

    })

    $(obj).addClass("currentItem");


    //ID重複チェック ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝
    var idArr = [];
    var duplicateIdArr = [];
    [].forEach.call(document.querySelectorAll('[id]'), function (elm) {
        var id = elm.getAttribute('id');
        if (idArr.indexOf(id) !== -1) {
            duplicateIdArr.push(id);
        } else {
            idArr.push(id);
        }
    });
    if (duplicateIdArr.length > 0) {
        console.log('IDの重複があります:', duplicateIdArr);
    } else {
        console.log('IDの重複はありません。');
    }
    //ID重複チェック ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝


} // =======================================



