package app

type Dbase struct {
	Data map[string]string
}

func NewDbase() Dbase {
	s := Dbase{
		Data: make(map[string]string),
	}
	return s
}

func (s *Dbase) SaveLink(url []byte, id string) {
	s.Data[id] = string(url)
}
