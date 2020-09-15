package main

import (
	"context"
	"github.com/A-ndrey/raspi-manage-bot/bot"
	"github.com/A-ndrey/raspi-manage-bot/configs"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := bot.Start(ctx, config)
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
