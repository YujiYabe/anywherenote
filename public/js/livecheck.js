
function exitApp() {
    window.close();
}

window.onload = function () {
    isEnableAppMode = true; // for debug
    waitSecondLiveCheck = 5; // for debug
    // 処理
    setInterval(function () { liveCheck() }, 1000);

};

function liveCheck() {
    // var dt = new Date();
    // var expireLiveTime = dt.setMinutes(dt.getSeconds() + 10);

    // var post_data = {
    //     'expireLiveTime': expireLiveTime
    // };

    $.ajax({
        // type: 'POST',
        // data: post_data,
        url: 'livecheck'

    })
        .then(
        // 通信成功
        function () {
            // $("#results").append(data);
            // console.log("ok");
        },
        // 通信失敗
        function () {
            if (isEnableAppMode) {
                window.open('about:blank', '_self').close();
            }
        });
}