package main

import (
	"flag"
	"log"

	"github.com/DevonTM/sankaku-viewer"
)

var (
	listen   = flag.String("l", ":8000", "Listen address")
	proxy    = flag.String("p", "", "Proxy address")
	username = flag.String("user", "", "Username")
	password = flag.String("pass", "", "Password")
)

func main() {
	flag.Parse()
	if *proxy != "" {
		if err := sankaku.UseProxy(*proxy); err != nil {
			log.Fatalln(err)
		}
	}

	if *username != "" && *password != "" {
		log.Println("Logging in as:", *username)
		if err := sankaku.Login(*username, *password); err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("Listening HTTP Server on:", *listen)
	if err := sankaku.ListenAndServe(*listen); err != nil {
		log.Fatalln(err)
	}
}
