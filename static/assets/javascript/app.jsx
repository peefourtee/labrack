import _ from "lodash";
import React from "react";

import store from "./store";
import {connectFeed} from "./actions";
import Device from "./Device.jsx";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {data: store.getState()};
  }

  componentDidMount() {
    this.unsubscribe = store.subscribe(() => {
      this.setState({data: store.getState()});
    });
    this.ws = connectFeed(store);
  }

  componentWillUnmount() {
    this.ws.close();
    this.unsubscribe();
  }

  renderDevice(device) {
    return <Device key={"device" + device.device} id={device.device} telemetry={device.telemetry} />;
  }

  renderConnected(connected) {
    if (connected) {
      return <div className="alert alert-success" role="alert">connected</div>;
    } else {
      return <div className="alert alert-danger" role="alert">disconnected</div>;
    }
  }

  renderErrors(errors) {
    if (errors && errors.length > 0) {
      return <div className="alert alert-danger" role="alert">{errors[errors.length-1]}</div>;
    }
  }

  render() {
    let data = this.state.data;

    let deviceIDs = _.keys(data.devices);
    let devices = [];
    for (let id of deviceIDs) {
      devices.push(data.devices[id]);
    }

    return (
      <div>
        {this.renderConnected(data.connected)}
        {this.renderErrors(data.errors)}
        {devices.map(this.renderDevice)}
      </div>
    );
  }
}
export default App;
