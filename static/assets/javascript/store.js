import _ from "lodash";
import {createStore} from "redux";

const TELEMETRY_HISTORY_LENGTH = 5;

const initialState = {
  devices: {},
  connected: false
};

const reduceAction = (state = initialState, action) => {
  switch (action.type) {
  case "TELEMETRY_RX": {
    state = _.cloneDeep(state);
    let data = action.data;

    if (!state.devices.hasOwnProperty(data.Device)) {
      state.devices[data.Device] = {device: data.Device, telemetry: {}};
    }

    let device = state.devices[data.Device];
    let telemetry = device.telemetry;

    if (!telemetry.hasOwnProperty(data.Name)) {
      telemetry[data.Name] = [];
    }
    telemetry[data.Name] = [data].concat(telemetry[data.Name]);

    if (telemetry[data.Name].length > TELEMETRY_HISTORY_LENGTH) {
      telemetry[data.Name] = telemetry[data.Name].slice(0, TELEMETRY_HISTORY_LENGTH);
    }
    device.telemetry = telemetry;
    return state;
  }
  case "SERVER_CONNECT": {
    state = Object.assign({}, state);
    state.connected = true;
    return state;
  }
  case "SERVER_DISCONNECT": {
    state = Object.assign({}, state);
    state.connected = false;
    return state;
  }
  default:
    return state;
  }
};

export default createStore(reduceAction);
