
function exitApp() {
    window.close();
}

window.onload = function () {
    isEnableAppMode = true; // サーバからのハートビート受け取り
    waitSecondLiveCheck = 5;  // ハートビート切断許容時間：秒
    // 処理
    setInterval(function () { liveCheck() }, 1000);

};

function liveCheck() {

    // var post_data = {
    //     'expireLiveTime': expireLiveTime
    // };

    if (isEnableAppMode) {

        $.ajax({
            url: 'livecheck',
            // type: 'POST',
            // data: post_data,

        })
            .then(
            // 通信成功
            function () {
                // $("#results").append(data);
                // console.log("ok");
            },
            // 通信失敗
            function () {
                window.open('about:blank', '_self').close();
            });
    }
}