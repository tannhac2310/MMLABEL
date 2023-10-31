var action = {
    //TODO: add action
    "action_ping": {
        "request": {"seq":1,"action":"ping","data":{"message":"ping"}},
        "response": {"seq":1,"action":"ping","data":{"message":"ping"}},
    },
    "action_broadcast": {
        "request": {"seq":1,"action":"broadcast","data":{"user_ids": ["string"], "message": "string"}},
    },
    "event_broadcast": {
        "response": {"event": "broadcast", "data": {"message":"string"}, "seq": 0},
        "description": "Receive when another client want to broadcast anything",
    },
    "event_iot_device_state_change": {
          "response": {"event": "iot_device_state_change", "data": {"pondId":"string", "iotDeviceId": "string", "alarmStatus": "OK/WARNING", "feedbackStart": true, "commandStart": true}, "seq": 0},
          "description": "Receive when iot device change alarm status",
    },
    "event_pond_mode_change": {
          "response": {"event": "pond_mode_change", "data": {"pondId":"string", "mode":"auto/manual"}, "seq": 0},
          "description": "Receive when pond change mode",
    },
    "event_edge_device_state_change": {
          "response": {"event": "edge_device_state_change", "data": {"pondId":"string", "pondCode":"string", "pondName": "string", "edgeDeviceStatus": "online/offline"}, "seq": 0},
          "description": "Receive when edge device change state",
    }
};

$( document ).ready(function() {
    $("#chat_input").hide();
    $('#server_url').val(location.hostname+ ":" + location.port + "/gezu/ws");
    $("#btn_connect").click(function () {
        let btn = $("#btn_connect");
        if (btn.text().trim() === "CONNECT") {
            connectSocket($('#server_url').val());
        } else {
            socket.close();
        }
    });

    $("#btn_send").click(function () {
        socket.send($("#json-request").text());
    });

    for (var k in action) {
        var elem = `<li class="mdl-list__item">
                <span class="mdl-list__item-primary-content">
                  ` + k + `
                </span>
            </li>`
        $(".demo-list-item").append(elem);
    }

    function getJson() {
        return action["action_ping"];
    }

    var editorRequest = new JsonEditor('#json-request', getJson()["request"]);
    var editorResponse = new JsonEditor('#json-response', getJson()["response"]);
    // editor.load(getJson());

    var itemSelected = null;
    $(".mdl-list__item").click(function (e) {
        if (itemSelected != null) {
            itemSelected.toggleClass("active");
        }

        itemSelected = $(this);
        itemSelected.toggleClass("active");

        var data = action[$(this).text().trim().toLowerCase()];
        editorRequest.load(data["request"] || "empty");
        editorResponse.load(data["response"] || "empty");


        $('#description').text( data["description"] || "");

        $('#btn_send').prop('disabled', data["request"] === undefined);
    })
});
