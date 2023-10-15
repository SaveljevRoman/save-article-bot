package files

import (
	"bot/lib/e"
	"bot/lib/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const (
	defaultPerm = 0774
)

func NewStorage(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(p *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can not save page", err) }()

	filePath := filepath.Join(s.basePath, p.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fileName, err := fileName(p)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(Username string) (p *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can not save page", err) }()

	filePath := filepath.Join(s.basePath, p.UserName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(filePath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can not remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can not remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can not check file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can not check if file %s exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	const ErrorDecodePage = "can not decode page"
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap(ErrorDecodePage, err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap(ErrorDecodePage, err)
	}

	return &p, nil
}
