package storage

// Dbase - структура базы данных
type Dbase struct {
	Data      map[string]string
	UserLinks map[string][]string
}

// NewDbase - создание новой базы данных
func NewDbase() Dbase {
	s := Dbase{
		Data:      make(map[string]string),
		UserLinks: make(map[string][]string),
	}
	return s
}

// SaveLink - сохранение ссылки
func (s *Dbase) SaveLink(url []byte, id string) {
	s.Data[id] = string(url)
}

// SaveUserLink - сохранение ссылки пользователя
func (s *Dbase) SaveUserLink(id string, link string) {
	s.UserLinks[id] = append(s.UserLinks[id], link)
}

// LoadLink - загрузка ссылки
func (s *Dbase) LoadUserLink(id string) ([]string, bool) {
	if s.UserLinks[id] == nil {
		return nil, false
	}
	return s.UserLinks[id], true
}
