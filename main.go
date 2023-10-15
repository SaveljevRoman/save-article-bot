package main

import (
	tgClient "bot/clients/telegram"
	"bot/consumer/event-consumer"
	"bot/events/telegram"
	"bot/lib/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.NewTgProcessor(
		tgClient.NewClient(tgBotHost, mustToken()),
		files.NewStorage(storagePath),
	)

	log.Printf("service start")

	if err := event_consumer.NewEConsumer(eventsProcessor, eventsProcessor, batchSize).Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("token-bot-token", "", "token for access to telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
