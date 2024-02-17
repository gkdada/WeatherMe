package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	HttpPort      int
	WeatherApiKey string
	WeatherUrl    string
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, errors.New("error opening config file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cnf := Config{
		//default value(s)
		HttpPort: 8192,
	}
	err = decoder.Decode(&cnf)
	if err != nil {
		return nil, errors.New("error decoding config file: " + err.Error())
	}

	if len(cnf.WeatherApiKey) == 0 {
		return nil, errors.New("'WeatherApiKey' is required for querying openweathermap")
	}

	//prepend a forward slash if it doesn't have it already
	if len(cnf.WeatherUrl) == 0 || cnf.WeatherUrl[0] != '/' {
		cnf.WeatherUrl = "/" + cnf.WeatherUrl
	}

	if cnf.HttpPort == 0 || cnf.HttpPort > 65534 {
		return nil, errors.New("invalid port number in configuration. Needs to be between 1 and 65534")
	}

	return &cnf, nil
}
