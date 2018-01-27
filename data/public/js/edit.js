$(function () {
    load_drop_zone();

    $("#flexbox_cent_line").hover(
        function (e) {
            // $(this).css({ background: "blue" });
            // $("#separate_line").css("width", "15px")
        },
        function (e) {
            // $("#separate_line").css("width","3px")
        }
    );

    $("#flexbox_cent_line").draggable({

        axis: "x",

        containment: $(".flexbox_container"),
        start: function (e, ui) {
            $(this).addClass('dragging');
            console.log("  zz");
        },

        drag: function (e, ui) {

            // divider-1
            if (ui.helper[0].id === "flexbox_cent_line") {

                // let 2 flow
                $("#flexbox_rght_pane").css("flex", "1");

                // let 1 move
                $("#flexbox_left_pane").css("flex", "0 1 " + (ui.offset.left) + "px");

            }
        },
        stop: function (event, ui) {
            $(this).removeClass('dragging');
        }

    });
});

function load_drop_zone() {

    // 保護機能によって更新されなかったファイルに対して、そのファイルの下にメッセージを出す。
    // アップロード処理エリア設定
    Dropzone.options.fileUpload = {
        // maxFilesize: json_setting.file_size_limite, // MB単位 設定の値をページ読み込み時に反映
        // acceptedFiles: json_setting.list_valid_extension,

        success: function (event, res, xhr) {

            $('#source_return_value').html(res);

            // remake_table();
            makePageList();

            if (res.code != 0) {

            } else if (!res.error) {

            } else {
                console.log('アップロード失敗');


            }
        },

        // dictDefaultMessage: json_dictionary.message_for_upload
    }
}



