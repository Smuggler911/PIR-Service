package main

import (
	"MakeWish-serverSide/config"
	"MakeWish-serverSide/internal/api"
	"MakeWish-serverSide/internal/server"
	"log"
)

func init() {
	config.ConnectDB()
}
func main() {
	conf, _ := config.LoadConfig()
	handlers := new(api.Handler)
	srv := new(server.Server)

	if err := srv.Run(conf.Port, handlers.InitRoutes()); err != nil {
		log.Fatalln("Error start server: " + err.Error())
	}
}
