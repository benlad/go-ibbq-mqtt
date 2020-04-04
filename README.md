# go-iBBQ MQTT Publisher Example
Inspired by sworisbreathing/go-ibbq, this is a simple app that connects to an iBBQ over BLE using sworisbreathing/go-ibbq. It publishes the data it receives to an MQTT channel using github.com/eclipse/paho.mqtt.golang 

## Building

### Linux

```bash
$ GOOS=linux go build
```

### OS X

```bash
$ GOOS=darwin go build
```

## Usage

```bash
$ LOGXI=* ./go-ibbq-mqtt
12:56:06.920508 INF main Connecting to device
12:56:13.419140 INF main Connected to device
12:56:13.433666 INF main Received battery data batteryPct: 96
12:56:14.123995 INF main Received temperature data temperatures: [19 18]
12:56:16.164030 INF main Received temperature data temperatures: [19 18]
12:56:18.503975 INF main Received temperature data temperatures: [19 18]
12:56:20.453983 INF main Received temperature data temperatures: [19 18]
12:56:22.404003 INF main Received temperature data temperatures: [19 18]
^C12:56:24.377496 INF main Disconnected # <- ctrl-C was pressed (SIGINT)
12:56:24.467517 INF main Exiting
$
```