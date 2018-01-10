
function exitApp() {
    window.close();
}

window.onload = function () {

    // 処理
    setInterval(function () { liveCheck() }, 1000);


};

function liveCheck() {


    $.ajax({
        url: "livecheck"
    })
        .then(
        // 通信成功
        function () {
            // $("#results").append(data);
            // console.log("ok");
        },
        // 通信失敗
        function () {
            console.log("ng");
            window.open('about:blank', '_self').close();
        });
}