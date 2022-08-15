package main

import (
	"ArvanWallet/services"
	"log"
	"os"
)

func main() {
	var e error
	e = services.ConnectDB()
	if e != nil {
		log.Print("Error in db connection")
	}

	e = services.ConnectRabbitMq()
	if e != nil {
		log.Print("Error in queue connection")
	}

	e = services.ConnectRedis()
	if e != nil {
		log.Print("Error in redis connection")
	}

	s := services.QueueService{
		Rabbit: services.RABBIT,
		Db:     services.DB,
		Redis:  services.REDIS,
		Ctx:    services.CTX,
	}
	go func() {
		s.ReceiveVoucherRegister("even")
	}()
	go func() {
		s.ReceiveVoucherRegister("odd")
	}()

	port := os.Args[1]
	GetRouter(port)
}
