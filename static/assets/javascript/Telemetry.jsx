import React from "react";

class TelemetryList extends React.Component {
  render() {
    if (!this.props.readings || this.props.readings.length <= 1) {
      return null;
    }
    return (
      <div>
        <label>Previous Readings</label>
        <ul>
          {this.props.readings.map(function(reading) {
            return <li key={"telemetry" + reading.Device + reading.Received}>{reading.Received}: {reading.Value}</li>;
          })}
        </ul>
      </div>
    );
  }
}

export default class Telemetry extends React.Component {
  render() {
    let {device, name, history} = this.props;
    return (
      <div>
        <h3>{name}: {history[0].Value}</h3>
        <TelemetryList device={device} readings={history} />
      </div>
    );
  }
}
