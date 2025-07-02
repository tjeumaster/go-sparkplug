package main

import (
	"time"

	"github.com/tjeumaster/go-sparkplug/sparkplug"
)

type Device struct {
	ID string
}

func (d *Device) GetDeviceID() string {
	return d.ID
}

func (d *Device) GetMetricValues() map[string]any {
	return map[string]any{
		"Temperature": 25.5,
		"Humidity":    60,
		"Status":      "OK",
	}
}

func (d *Device) GetMetricValueByName(name string) any {
	v, ok := d.GetMetricValues()[name]
	if !ok {
		return nil
	}

	return v
}

func main() { 
	c := &sparkplug.SparkplugConfig{
		Host:     "localhost",
		Port:     1883,
		Username: "username",
		Password: "password",
		ClientID: "dev",
		GroupID:  "RecipePlus",
		NodeID:   "PKV31-01",

	}

	client := sparkplug.NewClient(c)
	client.Start()

	d := &Device{ID: "Checkweigher1"}

	client.PublishDBIRTH(d)

	time.Sleep(2 * time.Second)

	client.PublishDDEATH(d)

	time.Sleep(2 * time.Second)

	client.PublishDBIRTH(d)
}