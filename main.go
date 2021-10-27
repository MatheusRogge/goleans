package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"com.github/MatheusRogge/budget-service/pkg/runtime/grain"
)

func main() {
	repository := grain.NewInMemoryRepository()

	grain, err := repository.GetGrain("Player", "Player1")
	if err != nil {
		panic(err)
	}

	grain.SetMessageHandler(func(msg interface{}) error {
		log.Printf("Received message %+v", msg)
		return nil
	})

	for i := 0; i < 200; i++ {
		grain.SendMessage(i)
	}

	signal_chan := make(chan os.Signal, 1)

	signal.Notify(
		signal_chan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-signal_chan
			switch s {

			case syscall.SIGINT:
				repository.Stop()
				exit_chan <- 0

			case syscall.SIGTERM:
				repository.Stop()
				exit_chan <- 0

			default:
				fmt.Println("Unknown signal.")
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan
	os.Exit(code)
}
