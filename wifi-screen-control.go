package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func IsThereAnyBodyOutThere(device string) bool {
	out, err := exec.Command("iw", "dev", device, "station", "dump").Output()
	if err != nil {
		log.Fatal(err)
	}
	devicePresent := strings.Contains(string(out), "Station ")
	return devicePresent
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

func Loop(device string, pollingInterval time.Duration) {
	// first, make sure the monitor is on on exit
	defer SwitchMonitor(true)

	// properly exit on signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Main loop
	// but first: we don't know in which state the monitor is, let's force it
	previousDevicePresent := IsThereAnyBodyOutThere(device)
	SwitchMonitor(previousDevicePresent)

	for {
		devicePresent := IsThereAnyBodyOutThere(device)

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
					Name:  "interval, n",
					Value: 10,
					Usage: "Polling interval for the wifi status, in seconds",
				},
			},
			Action: func(c *cli.Context) error {
				device := c.Args().Get(0)
				pollingInterval := time.Duration(c.Int("interval")) * time.Second
				fmt.Printf("Checking wifi AP device %v every %v\n", device, pollingInterval)

				Loop(device, pollingInterval)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
