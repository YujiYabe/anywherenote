"use strict";

function addNote(Obj) {


    var target_url = 'addnote';
    var post_data = {
        'new_note_name': $('#new_note_name').val(),
        'new_note_address': $('#new_note_address').text(),
    };

    $.ajax({
        type: 'POST',
        url: target_url,
        data: post_data
    })
        .then(
        // 1つめは通信成功時のコールバック
        function (data) {
            location.reload();
        },
        // 2つめは通信失敗時のコールバック
        function () {
            console.log("読み込み失敗");
        });
    ;


} // =======================================



function addPage(Obj) {
    // console.log($(Obj).attr('data-address'));
    var target_url = 'addpage';
    var post_data = {
        'note_address': $(Obj).attr('data-address'),
        'note_id': $(Obj).attr('data-note_id'),
    };

    $.ajax({
        type: 'POST',
        url: target_url,
        data: post_data
    })
        .then(
        // 1つめは通信成功時のコールバック
        function (data) {
            // location.reload();
            // console.log(data);
            $('#source_return_value').text(data);
            makePageList();
        },
        // 2つめは通信失敗時のコールバック
        function () {
            console.log("読み込み失敗");
        });
    ;
} // =======================================




function updateNote(Obj) {

    var target_url = 'updatenote';
    var post_data = {
        'new_note_name': $(Obj).next().children('input').val(),
        'note_id': $(Obj).next().children('input').attr('data-note_id'),
    };

    $.ajax({
        type: 'POST',
        url: target_url,
        data: post_data
    })
        .then(
        // 1つめは通信成功時のコールバック
        function (data) {
            location.reload();

        },
        // 2つめは通信失敗時のコールバック
        function () {
            console.log("読み込み失敗");
        });
    ;

} // =======================================



function updatePage() {

    var ed = tinyMCE.get('page_body');
    var data = ed.getContent();

    var target_url = 'updatepage';
    var post_data = {
        'note_id': $('#note_id').text(),
        'note_address': $('#note_address').text(),
        'page_id': $('#page_id').text(),
        'page_title': $('#page_title').val(),
        'page_body': data
    };

    $.ajax({
        type: 'POST',
        url: target_url,
        data: post_data,
    })
        .then(
        // 1つめは通信成功時のコールバック
        function (data) {
            // location.reload();
            $('#source_return_value').text(data);
            makePageList();

            $('.update_successed').show();
            $('.update_successed').fadeOut(3000);

        },
        // 2つめは通信失敗時のコールバック
        function () {
            console.log("読み込み失敗");
        });
    ;


} // =======================================





function deleteNote(Obj) {

    var target_url = 'deletenote';
    var post_data = {
        'note_id': $(Obj).prevAll().eq(1).children('input').attr('data-note_id'),
    };

    $.ajax({
        type: 'POST',
        url: target_url,
        data: post_data
    })
        .then(
        // 1つめは通信成功時のコールバック
        function (data) {
            location.reload();
        },
        // 2つめは通信失敗時のコールバック
        function () {
            console.log("読み込み失敗");
        });
    ;

} // =======================================

function deletePage() {


    var target_url = 'deletepage';
    var post_data = {
        'note_address': $('#note_address').text(),
        'page_id': $('#page_id').text(),
    };

    $.ajax({
        type: 'POST',
        url: target_url,
        data: post_data
    })
        .then(
        // 1つめは通信成功時のコールバック
        function (data) {
            location.reload();
        },
        // 2つめは通信失敗時のコールバック
        function () {
            console.log("読み込み失敗");
        });
    ;


} // =======================================

