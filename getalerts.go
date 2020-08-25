package main

import (
	"fmt"
	"log"
	"log/syslog"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/alerts-v2"
	"github.com/pcktdmp/cef/cefevent"
)

func main() {

	logwriter, err := syslog.New(syslog.LOG_NOTICE, "akamai")
        if err != nil {
		log.Fatal(err)
        }

	log.SetOutput(logwriter)

        config, err := edgegrid.Init("~/.edgerc", "papi")
        if err != nil {
		log.Fatal(err)
        }

	alerts.Init(config)

        response, err := alerts.ListActiveFirings()
        if err != nil {
		log.Fatal(err)
        }
	
	for _, firing := range response.Firings {
		var classid = "ACTIVE_ALERT"
		if !firing.EndTime.IsZero() {
			classid = "CLEARED ALERT"
		}

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
			DeviceEventClassId: classid,
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

}
