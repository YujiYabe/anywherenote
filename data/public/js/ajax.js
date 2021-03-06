"use strict";

function createNote(Obj) {

    var datetime = moment().format("YYYY/MM/DD HH:mm:ss");

    var target_url = 'createnote';
    var post_data = {
        'note_name': $('#new_note_name').val(),
        'note_address': $('#new_note_address').val().trim(),
        'note_update': datetime,
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



function createPage(Obj) {
    // console.log($(Obj).attr('data-address'));
    var datetime = moment().format("YYYY/MM/DD HH:mm:ss");

    var target_url = 'createpage';

    var post_data = {
        'note_address': $(Obj).attr('data-address'),
        'note_id': $(Obj).attr('data-select_note_id'),
        'page_update': datetime,
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
                // console.log(data);
                $('#source_dataset_list').text(data);
                makePageList();
            },
            // 2つめは通信失敗時のコールバック
            function () {
                console.log("読み込み失敗");
            });
    ;
} // =======================================




function updateNote(Obj) {

    // var note_star = $(Obj).nextAll().eq(0).children('[name=starint]').text();
    var note_name = $(Obj).nextAll().eq(1).children('input').val();
    var note_id = $(Obj).nextAll().eq(1).children('input').attr('data-note_id');
    var note_star = $(Obj).nextAll().eq(1).children('input').attr('data-star_int');
    var datetime = moment().format("YYYY/MM/DD HH:mm:ss");


    var target_url = 'updatenote';
    var post_data = {
        'note_id': note_id,
        'note_star': note_star,
        'note_name': note_name,
        'note_update': datetime,
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
                $(Obj).nextAll().eq(1).children('span').show();
                $(Obj).nextAll().eq(1).children('input').hide();
                $(Obj).nextAll().eq(1).children('span').fadeOut(3000,
                    function () { $(Obj).nextAll().eq(1).children('input').show() }

                );


            },
            // 2つめは通信失敗時のコールバック
            function () {
                console.log("読み込み失敗");
            });
    ;

} // =======================================



function updatePage() {

    // var ed = tinyMCE.get('page_body');
    // var data = ed.getContent();
    var datetime = moment().format("YYYY/MM/DD HH:mm:ss");

    var target_url = 'updatepage';
    var post_data = {
        'note_id': $('#note_id').text(),
        'note_address': $('#note_address').text(),

        'page_id': $('#page_id').text(),
        'page_star': $('#page_title').attr('data-star_int'),
        'page_title': $('#page_title').val(),
        'page_body': $('#page_body').html(),

        'page_update': datetime,

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
                $('#source_dataset_list').text(data);
                makePageList();

                $('.update_page_successed').show();
                $('.update_page_successed').fadeOut(3000);

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
        'note_id': $('#note_id').text(),
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
                // location.reload();

                $('#source_dataset_list').text(data);
                makePageList();

            },
            // 2つめは通信失敗時のコールバック
            function () {
                console.log("読み込み失敗");
            });
    ;


} // =======================================

