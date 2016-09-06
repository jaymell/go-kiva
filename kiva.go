package main

import (
	"os"
	"fmt"
	"github.com/rumdrums/go-kiva/kiva"
	"strconv"
)

func Client() *kiva.Client {
	var clientConfig kiva.Config
	return kiva.New(&clientConfig)
}

func main() {
	//PrintRawLoansJson()

 	cli := Client()
 	var loanIDs []int
    for _, v := range os.Args[1:] {
    	if id, e := strconv.Atoi(v); e == nil {
      	  loanIDs = append(loanIDs, id)
    	}
    	//loanIDs[i], _ = strconv.Atoi(v)
    }
    loans, err := cli.GetLoansByID(loanIDs...)
    if err != nil {
		fmt.Println(err)
	}
	fmt.Println("loans ", loans)

	for k, v := range loans {
		fmt.Println(k, v)
	}	

    for _, v := range loanIDs {
      lenders, err := cli.GetLoanLenders(v)
      if err != nil {
      	fmt.Println(err)
      	break
      } 
      for _, w := range lenders {
      	fmt.Println("lender: ", w)
      }
      similar, err := cli.GetSimilarLoans(v)
      if err != nil {
        fmt.Println(err)
        break
      } 
      for _, w := range similar {
        fmt.Println("similar: ", w)
      }
    }
}

