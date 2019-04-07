package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func IsThereAnyBodyOutThere(device string, threshold int) bool {
	out, err := exec.Command("iw", "dev", device, "station", "dump").Output()
	if err != nil {
		log.Fatal(err)
	}
	// out: either nothing or:
	// Station xx:xx:xx:xx:xx:xx (on wlan0)
	// 	signal:  	-42 dBm
	// get station "distance" by parsing its signal
	stationsSignal := make(map[string]int)
	var currentStation string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		switch fields := strings.Fields(scanner.Text()); fields[0] {
		case "Station":
			currentStation = fields[1]
		case "signal:":
			signal, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}
			stationsSignal[currentStation] = signal
			fmt.Printf("Found station %v with signal %v dBm\n", currentStation, signal)
		}
	}
	closestStationSignal := math.MinInt64
	for _, signal := range stationsSignal {
		if signal > closestStationSignal {
			closestStationSignal = signal
		}
	}
	// considere station present if signal is higher than threshold
	stationPresent := closestStationSignal != math.MinInt64 && closestStationSignal >= threshold
	return stationPresent
}

func SwitchMonitor(on bool) {
	var mode string
	if on {
		mode = "on"
	} else {
		mode = "off"
	}
	fmt.Printf("Switching monitor %v\n", mode)
	err := exec.Command("xset", "dpms", "force", mode).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Loop(device string, threshold int, pollingInterval time.Duration) {
	// first, make sure the monitor is on on exit
	defer SwitchMonitor(true)

	// properly exit on signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Main loop
	// but first: we don't know in which state the monitor is, let's force it
	previousDevicePresent := IsThereAnyBodyOutThere(device, threshold)
	SwitchMonitor(previousDevicePresent)

	for {
		devicePresent := IsThereAnyBodyOutThere(device, threshold)

		if devicePresent != previousDevicePresent {
			SwitchMonitor(devicePresent)
			previousDevicePresent = devicePresent
		}

		// loop or quit
		select {
		case <-quit:
			fmt.Println("\nReceived an interrupt, stopping...")
			return
		case <-time.After(pollingInterval):
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "wifi-screen-control"
	app.Usage = "Checking wifi AP for connected stations, controlling monitor on/off state."

	app.Commands = []cli.Command{
		{
			Name:      "watch",
			Aliases:   []string{"w"},
			Usage:     "watch wifi AP for stations to connect",
			ArgsUsage: "DEVICE",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "signal-threshold, s",
					Value: -100,
					Usage: "Signal threshold (in dBm) above it the station is considered present",
				},
				cli.IntFlag{
					Name:  "interval, n",
					Value: 10,
					Usage: "Polling interval for the wifi status, in seconds",
				},
			},
			Action: func(c *cli.Context) error {
				device := c.Args().Get(0)
				signalThreshold := c.Int("signal-threshold")
				pollingInterval := time.Duration(c.Int("interval")) * time.Second
				fmt.Printf("Checking wifi AP device %v (for stations signal>=%v) every %v\n", device, signalThreshold, pollingInterval)

				Loop(device, signalThreshold, pollingInterval)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
