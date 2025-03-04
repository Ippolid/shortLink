package storage

type Dbase struct {
	Data      map[string]string
	DataUsers map[string]string
}

func NewDbase() Dbase {
	s := Dbase{
		Data:      make(map[string]string),
		DataUsers: make(map[string]string),
	}
	return s
}

func (s *Dbase) SaveLink(url []byte, id string) {
	s.Data[id] = string(url)
}

func (s *Dbase) SaveUsersLink(userid string, id string) {
	s.DataUsers[id] = userid
}
