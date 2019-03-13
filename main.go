package main

import (
	"log"
	"spiderman-agent/app"
	"spiderman-agent/service"
)

func main() {

	//init server
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := app.GetApp().Prepare(); err != nil {
		panic(err)
	}
	//register service to etcd
	s, err := service.RegisterService()
	if err != nil {
		panic(err)
	}
	//monitor task change
	err = service.MonitorTask(s)
	if err != nil {
		panic(err)
	}
}
