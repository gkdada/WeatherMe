package main

import (
	"fmt"
	"sync"

	"github.com/gkdada/WeatherMe/config"
	"github.com/gkdada/WeatherMe/server"
)

func main() {
	cnf, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Error loading Config.", err)
		return
	}

	var wg sync.WaitGroup

	//we don't really need to have a goroutine for this. We might as well have the HttpServer running in main thread.
	//but this adds flexibility at a little cost and sets the platform for more threads/functionality as required.
	wg.Add(1)
	go func() {
		ws := server.NewServer(cnf)
		ws.HttpServer(&wg)
	}()

	wg.Wait()
}
