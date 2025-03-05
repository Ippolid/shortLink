package storage

type Dbase struct {
	Data      map[string]string // shortID -> originalURL
	DataUsers map[string]string // shortID -> userID
	Deleted   map[string]bool   // shortID -> true/false
}

func NewDbase() Dbase {
	return Dbase{
		Data:      make(map[string]string),
		DataUsers: make(map[string]string),
		Deleted:   make(map[string]bool),
	}
}

func (s *Dbase) SaveLink(url []byte, id string) {
	s.Data[id] = string(url)
}

func (s *Dbase) SaveUsersLink(userid string, id string) {
	s.DataUsers[id] = userid
}
