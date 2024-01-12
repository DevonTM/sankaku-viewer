package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/DevonTM/sankaku-viewer"
)

const VERSION = "v1.7.2"

var (
	listen   = flag.String("l", ":8000", "Listen address")
	proxy    = flag.String("p", "", "Proxy address")
	username = flag.String("user", "", "Username")
	password = flag.String("pass", "", "Password")
	cache    = flag.Bool("cache", true, "Enable cache for compressed files")
	version  = flag.Bool("version", false, "Print Version")
)

func main() {
	flag.Parse()
	if *version {
		fmt.Println("Sankaku Viewer " + VERSION)
		return
	}

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

	sankaku.CacheCompressed = *cache

	log.Println("Listening HTTP Server on:", *listen)
	if err := sankaku.ListenAndServe(*listen); err != nil {
		log.Fatalln(err)
	}
}
