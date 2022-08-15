package services

import (
	"ArvanWallet/requests"
	"ArvanWallet/resources"
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type QueueService struct {
	Rabbit *amqp.Connection
	Redis  *redis.Client
	Ctx    context.Context
	Db     *gorm.DB
}

//ReceiveVoucherRegister is consumer for voucher service
func (service QueueService) ReceiveVoucherRegister(routeKey string) {

	//create a rabbitMQ channel
	ch, e := service.Rabbit.Channel()
	if e != nil {
		return
	}
	defer ch.Close()

	//declare an exchange for the same topic in voucher service
	e = ch.ExchangeDeclare(
		"vouchers_topic", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if e != nil {
		return
	}

	q, e := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if e != nil {
		return
	}

	//bind queue due to even or odd routing key
	e = ch.QueueBind(
		q.Name,           // queue name
		"*."+routeKey,    // routing key
		"vouchers_topic", // exchange
		false,
		nil)
	if e != nil {
		return
	}

	//get data from queue
	msgs, e := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if e != nil {
		return
	}

	var forever chan struct{}

	//process the message fetched from queue
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			log.Printf(" [x] %s", d.RoutingKey)

			var receive resources.ReceiveVoucherResource
			if e := json.Unmarshal(d.Body, &receive); e == nil {
				mobileByte, _ := json.Marshal(receive.Data)
				mobile, _ := strconv.ParseInt(string(mobileByte), 10, 64)
				req := requests.AddTransactionRequest{
					Mobile: mobile,
					Amount: receive.Amount,
					Reason: "voucher_" + d.RoutingKey,
				}
				s := UserService{
					ServiceDB: service.Db,
					Redis:     service.Redis,
					Ctx:       service.Ctx,
				}
				//call AddTransaction service to update user status
				if e := s.AddTransaction(req); e == nil {
					d.Ack(false)
				}
			}
		}
	}()

	<-forever
}
