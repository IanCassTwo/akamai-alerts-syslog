package main

import (
	"fmt"
	"time"
	"log"
	"log/syslog"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/alerts-v2"
	"github.com/pcktdmp/cef/cefevent"
)

func main() {

	var lastRun map[string]alerts.Firing
	var thisRun map[string]alerts.Firing

	lastRun = make(map[string]alerts.Firing)
	thisRun = make(map[string]alerts.Firing)

	logwriter, err := syslog.New(syslog.LOG_NOTICE, "akamai")
        if err != nil {
		log.Fatal(err)
        }

	log.SetOutput(logwriter)

        config, err := edgegrid.Init("~/.edgerc", "alerts")
        if err != nil {
		log.Fatal(err)
        }

	alerts.Init(config)

	for range time.Tick(10 * time.Second) {

		response, err := alerts.ListActiveFirings()
		if err != nil {
			log.Fatal(err)
		}
	
		// Prepare
		thisRun = make(map[string]alerts.Firing)

		// New alerts
		for _, firing := range response.Firings {
			thisRun[firing.FiringID] = firing

			_, exists := lastRun[firing.FiringID] 
			if !exists {
				dolog("ACTIVE_ALERT", firing)
			}
		}

		// Cleared alerts
		for _, firing := range lastRun {
			_, exists := thisRun[firing.FiringID] 
			if !exists {
				dolog("CLEARED_ALERT", firing)
			}
		}

		// Save for next round
		lastRun = thisRun
	}
}

func dolog(classID string, firing alerts.Firing) {

	ext := make(map[string]string)
	ext["firingId"] = firing.FiringID
	ext["definitionId"] = firing.DefinitionID
	ext["startTime"] = firing.StartTime.String()
	ext["endTime"] = firing.EndTime.String()
	for k, v := range firing.FieldMap {
		ext[k] = fmt.Sprintf("%s", v);
	}
	
	event := cefevent.CefEvent{
		Version:            "0",
		DeviceVendor:       "Akamai",
		DeviceProduct:      firing.Service,
		DeviceVersion:      "1.0",
		DeviceEventClassId: classID,
		Name:               firing.Name,
		Severity:           "3",
		Extensions:         ext,
	}

	cef, err := event.Generate()
	if err != nil {
		log.Fatal(err)
	}

	log.Print(cef)
}
