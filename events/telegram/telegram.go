package telegram

import (
	"bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/storage"
	"context"
	"errors"
)

type TgProcessor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
	ctx     context.Context
}

type Meta struct {
	ChatId   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewTgProcessor(client *telegram.Client, storage storage.Storage) *TgProcessor {
	return &TgProcessor{
		tg:      client,
		storage: storage,
		ctx:     context.TODO(),
	}
}

func (tp *TgProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := tp.tg.Updates(tp.offset, limit)
	if err != nil {
		return nil, e.Wrap("can not get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	tp.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (tp *TgProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return tp.processMessage(event)
	default:
		return e.Wrap("can not process message", ErrUnknownEventType)
	}
}

func (tp *TgProcessor) processMessage(event events.Event) (err error) {
	defer func() { err = e.Wrap("can not process message", err) }()
	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err := tp.doCmd(event.Text, meta.ChatId, meta.UserName); err != nil {
		return err
	}

	return nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.ID,
			UserName: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can not get meta", ErrUnknownMetaType)
	}

	return res, nil
}
