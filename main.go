package main

import (
	"log"

	"github.com/pkg/browser"
)

func run() {
	configs := loadConfig(".")
	for _, conf := range configs {
		cs := CacheSystem{ListenAddr: conf.ListenAddr, DefaultServer: conf.DefaultServer, DefaultScheme: conf.DefaultScheme, CacheRoot: conf.CacheRoot, OfflineDomains: conf.OfflineDomains}
		go cs.Listen()
	}
	navigateAddr := "127.0.0.1:8200"
	go navigateServe(navigateAddr, configs)
	browser.OpenURL("http://" + navigateAddr)
}

func main() {
	//set log format show line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	run()
	select {}
}
