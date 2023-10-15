package telegram

import (
	"bot/lib/e"
	"bot/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (tp *TgProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return tp.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return tp.sendRandom(chatID, username)
	case HelpCmd:
		return tp.sendHelp(chatID)
	case StartCmd:
		return tp.SendHello(chatID)
	default:
		return tp.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (tp *TgProcessor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.Wrap("can not do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := tp.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExist {
		return tp.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := tp.storage.Save(page); err != nil {
		return err
	}

	if err := tp.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (tp *TgProcessor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can not command: can not send random", err) }()

	page, err := tp.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return tp.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := tp.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return tp.storage.Remove(page)
}

func (tp *TgProcessor) sendHelp(chatID int) error {
	return tp.tg.SendMessage(chatID, msgHelp)
}

func (tp *TgProcessor) SendHello(chatID int) error {
	return tp.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
