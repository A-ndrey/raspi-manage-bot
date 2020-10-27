package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/A-ndrey/raspi-manage-bot/bot"
	"github.com/A-ndrey/raspi-manage-bot/configs"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"github.com/A-ndrey/raspi-manage-bot/stats"
)

func main() {
	log.Println("Start application")
	defer log.Println("Shutdown application")

	config := configs.LoadConfig()

	if err := db.Init(); err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	monitoring := make(chan db.Measurement, 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		stats.StartMeasuring(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		stats.StartAnalyzing(ctx, config, monitoring)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := bot.Start(ctx, config, monitoring)
		if err != nil {
			log.Println(err)
			cancelFunc()
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancelFunc()
	}()

	wg.Wait()
}
