var socket = null;

var log = {
    info: function (content) {
        $(".content-log").append("<span class='mdl-color-text--green-800'>" + new Date().toJSON() + " " + content + "</span><br/>");
        var objDiv = document.getElementById("content-log");
        objDiv.scrollTop = objDiv.scrollHeight;
    },
    error: function (content) {
        $(".content-log").append("<span class='mdl-color-text--red-800'>" + new Date().toJSON() + " " + content + "</span><br/>");
        var objDiv = document.getElementById("content-log");
        objDiv.scrollTop = objDiv.scrollHeight;
    },
};


var connectSocket = function (url) {
    let schema = "ws://";
    if (location.protocol === 'https:') {
        schema = "wss://";
    }

    var token = $("#token").val();
    if (token === "") {
        log.error("Missing auth token");
        return;
    }
    url = schema + url+ "?token=" + token;
    server_url = url.substring(5);
    log.info("WS: connect to " + url);

    // getToken($("#userId").val());
    socket = new WebSocket(url);
    // socket.binaryType = 'arraybuffer';
    socket.onopen = function (e) {
        log.info("WS: connect SUCCESS");
        $('#btn_connect').text("DISCONNECT");
    };

    socket.onmessage = function (e) {
        receiveMessage(e.data)
    }

    socket.onerror = function (e) {
        console.error(e);
        log.error("WS: connect failed");
    };
    socket.onclose = function (e) {
        socket = null;
        log.error("WS: disconnect");
        $('#btn_connect').text("CONNECT");
    }
}

function receiveMessage(data) {
    if (JSON.parse(data).status === 'FAIL') {
        log.error(data);
        return
    }
    log.info(data);
}
