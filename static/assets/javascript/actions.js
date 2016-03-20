export const Actions = {
  telemetryRx: data => ({type: "TELEMETRY_RX", data: data}),
  serverConnect: () => ({type: "SERVER_CONNECT"}),
  serverDisconnect: () => ({type: "SERVER_DISCONNECT"})
};

export const connectFeed = (store) => {
  var ws = new WebSocket("ws://" + location.host + "/ws");
  ws.onerror = (err) => {
    throw new Error("error communicating with server: ", err);
  };

  ws.onopen = () => {
    store.dispatch(Actions.serverConnect());
  };

  ws.onmessage = (e) => {
    let message = JSON.parse(e.data);

    if (message.Type == "telemetry") {
      store.dispatch(Actions.telemetryRx(message.Data));
    } else {
      throw new Error("unknown websocket message type: ", message);
    }
  };

  ws.onclose = () => {
    store.dispatch(Actions.serverDisconnect());
    setTimeout(()=> {
      connectFeed(store);
    }, 5000);
  };

  return ws;
};
