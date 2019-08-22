package main

import (
	"fmt"
	"os"

	weather "github.com/vkorn/go-yahoo-weather"
)

const (
	appID        = "your-app-id"
	clientID     = "your-consumer-key"
	clientSecret = "your-consumer-secret"
)

func main() {
	yw := weather.NewProvider(appID, clientID, clientSecret)
	data, err := yw.Query("Oakland, CA", weather.Imperial)
	if err != nil {
		fmt.Printf("Got errror: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Temperature: %d \nCondition: %s", data.Observation.Condition.Temperature,
		data.Observation.Condition.Text)
}
