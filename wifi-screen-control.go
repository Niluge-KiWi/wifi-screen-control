package main

import "fmt"
import "time"
import "os/exec"
import "log"
import "strings"
import "os"
import "os/signal"
import "syscall"

const Device = "wlp0s29f7u2"
const PollingTime = 1 // in seconds

func IsThereAnyBodyOutThere() bool {
	out, err := exec.Command("iw", "dev", Device, "station", "dump").Output()
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

func main() {
	fmt.Printf("Checking wifi AP (device %v) for connected stations, controlling monitor on/off state.\n", Device)

	// switch monitor on on exit
	defer SwitchMonitor(true)

	// Main loop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	for {
		fmt.Print("Just checking... ")
		devicePresent := IsThereAnyBodyOutThere()
		fmt.Printf("%v\n", devicePresent)
		SwitchMonitor(devicePresent)

		// loop or quit
		select {
		case <-quit:
			fmt.Println("\nReceived an interrupt, stopping...")
			return
		case <-time.After(PollingTime * time.Second):
		}
	}
}
