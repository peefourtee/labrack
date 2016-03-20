package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/peefourtee/labrack"
	"github.com/peefourtee/labrack/app"
)

var (
	mockDevices = flag.Int("mock-devices", 0, "number of devices to mock telemetry data for")
	listenAddr  = flag.String("listen", ":8000", "address:port for the web server to listen on")
	sample      = flag.Duration("sample", 600*time.Millisecond, "sample inteveral")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: labrack [options] address1 address2...

Collect voltage and current stats for the given i2c addresses.

Options:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	if flag.NArg() == 0 && *mockDevices == 0 {
		fmt.Fprintf(os.Stderr, "no i2c addresses or number of mock-devices specified")
		flag.Usage()
	}

	telemetry := make(chan labrack.Telemetry, 100)
	errs := make(chan error)

	if *mockDevices > 0 {
		labrack.MockSource(telemetry, *mockDevices, *sample)
	}

	go setupI2CPolling(telemetry, errs)

	webserver(telemetry)
}

func setupI2CPolling(t chan<- labrack.Telemetry, errs chan<- error) {
	for _, s := range flag.Args() {
		// addresses are uint16 but packages we're using use int
		if i, err := strconv.ParseInt(s, 10, 17); err != nil {
			fmt.Fprintf(os.Stderr, "invalid address %q: %s", s, err)
			os.Exit(1)
		} else {
			go labrack.I2CSource(t, errs, int(i), *sample)
		}
	}
}

func webserver(t <-chan labrack.Telemetry) {
	mux := app.Mux(t)

	log.Printf("starting web server on %s", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, mux); err != nil {
		log.Fatal("failed to start web server: ", err)
		os.Exit(1)
	}
}
