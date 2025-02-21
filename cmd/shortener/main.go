package main

import (
	"github.com/Ippolid/shortLink/config"
	"github.com/Ippolid/shortLink/internal/app"
	"github.com/Ippolid/shortLink/internal/logger"
	"os"
)

func main() {
	host, adr := config.ParseFlags()
	//fmt.Println(host, adr)

	if envRunHost := os.Getenv("SERVER_ADDRESS"); envRunHost != "" {
		host = envRunHost
	}

	if envRunAdr := os.Getenv("BASE_URL"); envRunAdr != "" {
		adr = envRunAdr
		adr += "/"
	}

	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	k := app.NewDbase()
	s := app.New(&k, adr, host)
	if err := s.Start(); err != nil {
		panic(err)
	}
}
