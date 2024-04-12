package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/gkdada/WeatherMe/config"
)

type zipLookup struct {
	Zip     string  `json:"zip"`
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

type weatherLookup struct {
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type weatherProcessed struct {
	Condition string `json:"condition"`
	Feels     string `json:"feels"`
}

type WeatherServer struct {
	cnfg *config.Config
}

func NewServer(cnf *config.Config) *WeatherServer {
	return &WeatherServer{
		cnfg: cnf,
	}
}

func (ws *WeatherServer) roundFlt(val float64) float64 {
	return math.Round(val*10) / 10
}

func (ws *WeatherServer) processWeather(wl *weatherLookup) weatherProcessed {
	wlr := weatherProcessed{}
	if len(wl.Weather) > 0 {
		wlr.Condition = wl.Weather[0].Main
	}
	switch {
	case wl.Main.Temp > 301.48: //about 83 F
		wlr.Feels = "Hot"
	case wl.Main.Temp > 288.7: //about 60
		wlr.Feels = "Moderate"
	default:
		wlr.Feels = "Cold"
	}
	return wlr
}

func (ws *WeatherServer) validateLatLong(lat, long string) (float64, float64, error) {
	//we can use lat/long directly.
	//validate it.
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return 0, 0, errors.New("400 error decoding lattitude")
	}
	if latF < -90 || latF > 90 {
		return 0, 0, errors.New("400 lattitude should be between -90 and 90 degrees")
	}
	longF, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return 0, 0, errors.New("400 error decoding longitude")
	}
	if longF < -180 || longF > 180 {
		return 0, 0, errors.New("400 longitude should be between -180 and 180 degrees")
	}
	return latF, longF, nil
}

func (ws *WeatherServer) getWeatherForLatLong(latF, longF float64) (*weatherProcessed, error) {
	reqUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", latF, longF, ws.cnfg.WeatherApiKey)
	fmt.Println(reqUrl)
	res, err := http.Get(reqUrl)
	if err != nil {
		errStr := fmt.Sprintf("error '%s' querying openweathermap", err)
		return nil, errors.New(errStr)
	}
	if res.StatusCode != 200 {
		errStr := fmt.Sprintf("error %d querying openweathermap", res.StatusCode)
		return nil, errors.New(errStr)
	}
	bd, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("500 error reading WeatherApi output")
	}
	var wl weatherLookup
	err = json.Unmarshal(bd, &wl)
	if err != nil {
		return nil, errors.New("500 error decoding WeatherApi output")
	}
	wlr := ws.processWeather(&wl)
	return &wlr, nil
}

func (ws *WeatherServer) weatherForLatLong(w http.ResponseWriter, r *http.Request) {
	//3 possible parameters.
	lat := r.URL.Query().Get("lat")
	long := r.URL.Query().Get("long")
	var latF, longF float64
	var err error
	if len(lat) != 0 && len(long) != 0 {
		latF, longF, err = ws.validateLatLong(lat, long)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		//return an error
		http.Error(w, "400 Insufficient data. Need either lat & long or zip & [country]", http.StatusBadRequest)
		return
	}

	wlr, err := ws.getWeatherForLatLong(latF, longF)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(wlr)
}

func (ws *WeatherServer) HttpServer(wg *sync.WaitGroup) {
	defer wg.Done()

	http.HandleFunc(ws.cnfg.WeatherUrl, ws.weatherForLatLong)

	err := http.ListenAndServe(fmt.Sprintf(":%d", ws.cnfg.HttpPort), nil)
	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
