package main

import (
	lolzpayment "github.com/nlypage/lolz-payment/lolz-payment"
	"log"
	"sync"
)

func main() {
	lolzClient := lolzpayment.NewClient(lolzpayment.Options{
		Token: "",
	})

	var wg sync.WaitGroup
	wg.Add(1)

	lolzClient.HandlePayments(func(payment lolzpayment.Payment) error {
		log.Println(payment.OperationID)
		return nil
	}, nil)
	wg.Wait()
}
