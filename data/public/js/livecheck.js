"use strict";



function exitApp() {
    window.close();
}

window.onload = function () {

    var json_config_value = JSON.parse($('#source_user_config').text());

    var isEnableAppMode = json_config_value["IsEnableAppMode"]; // サーバからのハートビート受け取り
    var waitSecondLiveCheck = json_config_value["WaitSecondLiveCheck"] ;  // ハートビート切断許容時間：秒
    var WaitSecondInterval = json_config_value["WaitSecondInterval"] * 1000;  // ハートビート切断許容時間：秒

    // 処理
    if (isEnableAppMode) {
        setInterval(function () { liveCheck() }, WaitSecondInterval);
    }

};

function liveCheck() {
    // console.log($('#source_user_config').text());

    // var post_data = {
    //     'expireLiveTime': expireLiveTime
    // };


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