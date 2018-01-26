// File APIに対応していない場合はエリアを隠す
if (!window.File) {
    document.getElementById('image_upload_section').style.display = "none";
}

// ブラウザ上でファイルを展開する挙動を抑止
function onDragOver(event) {
    event.preventDefault();
}

//https://qiita.com/volpe28v/items/725e2b8c6a94ade505db

// Drop領域にドロップした際のファイルのプロパティ情報読み取り処理
function onDrop(event) {
    // ブラウザ上でファイルを展開する挙動を抑止
    event.preventDefault();

    // ドロップされたファイルのfilesプロパティを参照
    var files = event.dataTransfer.files;
    for (var i = 0; i < files.length; i++) {
        // ★最初の一件目で現在のページ内容を更新
        if (i == 0) {
            updatePage()
        }
        // 一件ずつアップロード
        imageFileUpload(files[i]);

        // 最後のアップロードでページを更新
        if (i == files.length - 1) {
            updatePage()
        }


    }
}

// ファイルアップロード
function imageFileUpload(f) {
    var formData = new FormData();
    formData.append('file', f);
    formData.append('note_address', $('#note_address').text());
    formData.append('page_id', $('#page_id').text());
    $.ajax({
        type: 'POST',
        contentType: false,
        processData: false,
        url: '/uploadfile',
        data: formData,
        dataType: 'json',
        success: function (data) {
            // メッセージ出したり、DOM構築したり。
        }
    });
}
