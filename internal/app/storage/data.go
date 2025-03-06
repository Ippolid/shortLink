package storage

import (
	"fmt"
	"github.com/Ippolid/shortLink/internal/app"
)

type Dbase struct {
	Data  map[string]string
	Users map[string][]string
}

func NewDbase() Dbase {
	s := Dbase{
		Data:  make(map[string]string),
		Users: make(map[string][]string),
	}
	return s
}

func (s *Dbase) SaveLink(url string, user string) (string, error) {
	id := app.GenerateShortID(url, user)
	if _, exist := s.Data[id]; exist {
		return "", fmt.Errorf(`id "%s" already exists`, id)
	}
	s.Data[id] = url
	s.Users[user] = append(s.Users[user], url)
	return id, nil
}

func (s *Dbase) GetLink(id string) (string, error) {
	url, exist := s.Data[id]
	if !exist {
		return "", fmt.Errorf(`id "%s" not found`, id)
	}
	return url, nil
}
