package main

import (
	"arcticwolf.com/cutler/dao"
	"arcticwolf.com/cutler/webserver"
)

func main() {
	backend := &dao.LocalCache{}
	webServer := webserver.New(backend)
	panic(webServer.ListenAndServe())
}
