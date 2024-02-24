package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/DevonTM/sankaku-viewer"
)

const VERSION = "v1.8.2"

var (
	listen        = flag.String("l", ":8000", "HTTP server listen address")
	proxy         = flag.String("p", "", "Proxy address,\nonly http and socks5 proxy are supported")
	username      = flag.String("user", "", "Sankaku username")
	password      = flag.String("pass", "", "Sankaku password")
	cache         = flag.Bool("cache", true, "Enable cache for compressed static files")
	CacheDuration = flag.Int("cache-duration", 600, "Set cache duration in seconds for cached posts data,\nset 0 to never expire, should be less than 1 hour")
	version       = flag.Bool("version", false, "Print Version")
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
	sankaku.SetCacheDuration(*CacheDuration)

	log.Println("Listening HTTP Server on:", *listen)
	if err := sankaku.ListenAndServe(*listen); err != nil {
		log.Fatalln(err)
	}
}
