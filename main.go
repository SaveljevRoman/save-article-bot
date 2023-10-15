package main

import (
	tgClient "bot/clients/telegram"
	"bot/consumer/event-consumer"
	"bot/events/telegram"
	"bot/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	fileStoragePath   = "file_storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	//file storage
	//s := files.New(fileStoragePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can not connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can not init storage: ", err)
	}

	eventsProcessor := telegram.NewTgProcessor(tgClient.NewClient(tgBotHost, mustToken()), s)

	log.Printf("service start")

	if err := event_consumer.NewEConsumer(eventsProcessor, eventsProcessor, batchSize).Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
