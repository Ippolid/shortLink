package storage

type Dbase struct {
	Data      map[string]string
	UserLinks map[string][]string
}

func NewDbase() Dbase {
	s := Dbase{
		Data:      make(map[string]string),
		UserLinks: make(map[string][]string),
	}
	return s
}

func (s *Dbase) SaveLink(url []byte, id string) {
	s.Data[id] = string(url)
}
func (s *Dbase) SaveUserLink(id string, link string) {
	s.UserLinks[id] = append(s.UserLinks[id], link)
}

func (s *Dbase) LoadUserLink(id string) ([]string, bool) {
	if s.UserLinks[id] == nil {
		return nil, false
	}
	return s.UserLinks[id], true
}
