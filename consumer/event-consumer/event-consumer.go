package event_consumer

import (
	"bot/events"
	"log"
	"time"
)

type EConsumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewEConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) EConsumer {
	return EConsumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (ec EConsumer) Start() error {
	for {
		gotEvents, err := ec.fetcher.Fetch(ec.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(time.Second)
			continue
		}

		for _, event := range gotEvents {
			log.Printf("got new event: %s", event.Text)

			if err := ec.processor.Process(event); err != nil {
				log.Printf("can not handle event: %s", err.Error())
				continue
			}
		}
	}
}
