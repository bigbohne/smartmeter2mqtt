package main

import (
	"log"
	"time"

	"github.com/alecthomas/kingpin/v2"
)

var (
	metertype = kingpin.Flag("type", "Type of Smartmeter [\"modbustcp\", \"iec62056\"]").Required().Envar("GOMETER_TYPE").String()
	metername = kingpin.Flag("name", "Name of Smartmeter").Required().Envar("GOMETER_NAME").String()
	interval  = kingpin.Flag("interval", "Interval of measurements in seconds").Default("15").Int()
	mqtturl   = kingpin.Flag("mqtt", "URL of MQTT Server").Envar("GOMETER_MQTTURL").String()
	urlinput  = kingpin.Flag("url", "Address of the Smartmeter").Envar("GOMETER_URL").String()
)

func main() {
	kingpin.Parse()

	// Create initial test measurement
	measurement, errMeasurement := createMeasurement(MeasurementSettings{
		metertype: *metertype,
		url:       *urlinput})
	if errMeasurement != nil {
		log.Fatalln(errMeasurement)
	}

	log.Println(measurement)

	var ticker = time.NewTicker(time.Second * time.Duration(*interval))

	var mqttclient *MQTTClient
	if len(*mqtturl) != 0 {
		var err error
		mqttclient, err = CreateMQTTClient(MQTTClientParams{
			name: *metername,
			url:  *mqtturl,
		})

		if err != nil {
			log.Fatalln(err)
		}
	}

	for {
		<-ticker.C

		var measurement, err = createMeasurement(MeasurementSettings{
			metertype: *metertype,
			url:       *urlinput})

		if err != nil {
			log.Fatal(err)
		}

		log.Println(measurement)

		if mqttclient != nil {
			mqttclient.Publish(measurement)
		}
	}

}
