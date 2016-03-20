import _ from "lodash";
import React from "react";

import Telemetry from "./Telemetry.jsx";

export default class Device extends React.Component {
  render() {
    let {id, telemetry} = this.props;

    // get the names/keys of all telemetry readings
    let readingNames = _.keys(telemetry);
    readingNames.sort();

    let readings = [];
    for (let name of readingNames) {
      readings.push(
        <Telemetry key={id+name} device={id} name={name} history={telemetry[name]} />
      );
    }
    return (
      <div>
        <h2>Address {this.props.id}</h2>
        {readings}
      </div>
    );
  }
}
Device.propTypes = {
  id: React.PropTypes.number,
  telemetry: React.PropTypes.object.isRequired
};
