package config

import "flag"

func ParseFlags() (string, string) {
	host := flag.String("a", "localhost:8080", "host")
	adr := flag.String("b", "http://localhost:8080/", "adress")
	flag.Parse()
	return *host, *adr
}
