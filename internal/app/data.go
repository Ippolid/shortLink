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

func (d *Dbase) SaveLink(url []byte, id string) {
	d.Data[id] = string(url)
}
