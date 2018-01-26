$(function () {

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
            console.log("");
        },

        drag: function (e, ui) {

            // divider-1
            if (ui.helper[0].id === "flexbox_cent_line") {

                // let 2 flow
                $(".flexbox_rght_pane").css("flex", "1");

                // let 1 move
                $(".flexbox_left_pane").css("flex", "0 1 " + (ui.offset.left) + "px");

                // divider-2
            }
        },
        stop: function (event, ui) {
            $(this).removeClass('dragging');
        }

    });
});
