### IoT Socket Streamer

Three different components that work together to stream sensor data from a Waveshare BME680 sensor to a web page with real-time updates via websockets. The components are:

- An `ESP32 Mini-1` microcontroller that uses [MicroPython](https://micropython.org) to read sensor data from a [Waveshare BME680 environmental sensor](https://www.waveshare.com/wiki/Bme680#:~:text=Overview.%20A%20tiny%20sensor%20breakout%20with%20bme680,also%20is%20compatible%20with%203.3V/5V%20voltage%20levels.). Wi-Fi wasn't an option here (5GHz only network), so the microcontroller sends the data out via a serial connection (UART) instead. All of the hard work is done with a micropython driver for the BME680 sensor located here [Micropython Driver for a BME680 breakout](https://github.com/adafruit/Adafruit_BME680) and here [Micropython Driver for a BME680](https://github.com/robert-hh/BME680-Micropython) The driver uses the I2C interface. On the ESP32 Mini-1, you can use PIN's 21 and 22 for the sensor.
  
- A client-side script written in [Go](https://go.dev) that reads from the serial port (the ESP32 Mini-1) and sends the data (via POST request in binary format) to a Golang server on a live domain (or localhost). The script uses the wonderful [serial](https://github.com/bugst/go-serial) package to read from the serial port. The [serial package](https://pkg.go.dev/go.bug.st/serial) documentation is very helpful and the package is *very easy* to use.

- A server-side script written in [Go](https://go.dev) that listens for the incoming connentions (sensor data), processes via a `POST` request handler, and sends it to a client (an index.html file) via [Gorilla WebSocket](https://github.com/gorilla/websocket). The server can be hosted on a live domain. The client is a simple HTML file (JS/HTML/Go templates) that connects to the server via a WebSocket connection and shows the sensor data updated in real-time.

---

### Usage

1. Connect the ESP32 to the BME680 sensor via I2C and connect the ESP32 to your computer via USB.

2. Start the server: `go run main.go`

3. Start the client script: `go run main.go`

4. Open the client in your browser: `http://localhost:8080`


---

### To Do

A lot that probably will never get finished 😅...most important thing is to get web sockets working on a live domain and not just locally, at moment you have to refresh the page for new sensor data...running it locally web sockets works...etc...refactor all of the code, add better error handling, make the landing page look nicer, add more sensors, etc...It's embarassingly basic...Go is so flexible, I had no idea it was as easy as it was for writing HTTP-based applications.
