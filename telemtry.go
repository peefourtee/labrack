package labrack

import (
	"log"
	"math/rand"
	"time"
)

type Telemetry struct {
	Device   int
	Name     string
	Value    float32
	Received time.Time
}

// for a given float v, return a new float that's a small delta alway from the
// current value, no larger than max.
func randFloat(v float32, max float32) float32 {
	tmp := rand.Float32()
	if rand.Intn(2) > 0 {
		tmp *= -1
	}

	for {
		if v+tmp < max && v+tmp >= 0 {
			break
		}
		if v+tmp > max {
			tmp -= rand.Float32()
		} else if v+tmp < 0 {
			tmp += rand.Float32()
		}
	}
	return v + tmp
}

func MockSource(c chan<- Telemetry, numDevices int, sample time.Duration) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numDevices; i++ {
		go func(id int) {
			voltage := Telemetry{
				Device: id,
				Name:   "voltage",
				Value:  rand.Float32() + float32(rand.Intn(30)),
			}
			current := Telemetry{
				Device: id,
				Name:   "current",
				Value:  rand.Float32() + float32(rand.Intn(5)),
			}

			send := func(t Telemetry) {
				select {
				case c <- t:
				default:
					log.Printf("failed to write mock telemetry for device %d", id)
				}
			}

			for {
				voltage.Value = randFloat(voltage.Value, 30)
				voltage.Received = time.Now().UTC()
				send(voltage)

				current.Value = randFloat(current.Value, 5)
				current.Received = time.Now().UTC()
				send(current)

				<-time.After(sample)
			}
		}(i)
	}
}
