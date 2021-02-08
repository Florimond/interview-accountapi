package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Florimond/interview-accountapi/service"
	"github.com/Florimond/interview-accountapi/service/contracts/account"
)

func main() {

	fmt.Println("Hello world")

	c := service.NewClient("http://localhost:8080/v1/", time.Minute)
	var acc account.Account

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	go func() {
		defer cancel()
		resp, err := c.FindByID(ctx, account.Provider, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
		if err != nil {
			log.Println(err.Error())
			return
		}
		resp.As(&acc)
		fmt.Println(acc)
	}()
	//cancel()
	select {
	case <-ctx.Done():
		//fmt.Println(ctx.Err())
	}

	/*
		acc := account.Account{
			contracts.RecordInfo{ID: "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"},
		}
		c.Delete(&acc)
		var accounts account.AccountList
		c.List(&accounts)
		fmt.Println(accounts)
	*/
}
