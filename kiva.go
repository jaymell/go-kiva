package main

import (
	"fmt"
	"github.com/rumdrums/go-kiva/kiva"
)

func Client() *kiva.Client {
	var clientConfig kiva.Config
	return kiva.New(&clientConfig)
}

func main() {
	//PrintRawLoansJson()

 	cli := Client()
	loans, err :=  cli.GetLoansById(1137156, 1128815)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("loans ", loans)
	for k, v := range loans {
		fmt.Println(k, v)
	}
}
