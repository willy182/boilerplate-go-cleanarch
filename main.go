package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/willy182/boilerplate-go-cleanarch/config"
	"github.com/willy182/boilerplate-go-cleanarch/utils"

	"github.com/joho/godotenv"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			utils.Log(log.ErrorLevel, fmt.Sprint(r), "main()", "recover_main")
		}
	}()

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env is not loaded properly")
		panic(err)
	}

	service := InitHSIService(config.Load())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	// pathSchema := "schemas"
	// jsonschema.Load(pathSchema)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		service.Serve()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case s := <-signals:
				fmt.Println(s.String())
				os.Exit(1)
			}
		}
	}()

	wg.Wait()
}
