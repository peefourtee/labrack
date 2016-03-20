package labrack

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hybridgroup/gobot/platforms/raspi"
)

type I2CError struct {
	Err     error
	Address int
}

func (e *I2CError) Error() string {
	return fmt.Sprintf("address %s: %s", e.Address, e.Err)
}

type i2cTelemetry struct {
	Voltage float32
	Amps    float32
}

func (t *i2cTelemetry) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return errors.New("invalid data length")
	}

	t.Voltage = math.Float32frombits(binary.LittleEndian.Uint32(data[0:4]))
	t.Amps = math.Float32frombits(binary.LittleEndian.Uint32(data[4:8]))
	return nil
}

func I2CSource(out chan<- Telemetry, errs chan<- error, address int, sample time.Duration) {
	r := raspi.NewRaspiAdaptor("i2c-" + strconv.Itoa(address))
	if err := r.I2cStart(address); err != nil {
		panic(fmt.Errorf("cannot start i2c comms for device %d", address))
	}
	for {
		data := i2cTelemetry{}

		if b, err := r.I2cRead(address, 8); err != nil {
			errs <- &I2CError{Address: address, Err: err}
		} else if len(b) != 8 {
			errs <- &I2CError{Address: address, Err: errors.New("invalid i2c response: expected 8 bytes")}
		} else if data.UnmarshalBinary(b); err != nil {
			errs <- &I2CError{Address: address, Err: errors.New("failed to decode i2c message")}
		} else {
			for _, t := range []Telemetry{
				{Device: address, Name: "voltage", Value: data.Voltage},
				{Device: address, Name: "current", Value: data.Amps},
			} {
				select {
				case out <- t:
				default:
					errs <- &I2CError{Err: errors.New("failed to send telemetry: channel full"), Address: address}
				}
			}
		}
		<-time.After(sample)
	}
}
