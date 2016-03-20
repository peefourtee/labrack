# labrack

There's an arduino talking to a couple different devices to monitor voltage and
current of various electronics in the lab rack.  labrack is a webapp serving
as a single pane of glass for all of these readings.

* telemetry data is polled from the arduino via i2c
* telemetry data is pushed to the client app via websockets

# dev

* run server
  ```
  docker build -t labrack:latest .
  docker run --rm -it --name labrack-running -p 8000:8000 -v //c/Users/foo/src/go/src/github.com/peefourtee/labrack:/go/src/github.com/peefourtee/labrack
  labrack:latest /bin/bash
  ```

* start web build pipeline
  ```
  cd static/
  npm install
  npm run dev
  ```
