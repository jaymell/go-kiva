package main

import (
	"fmt"
	"github.com/jaymell/go-kiva/kiva"
	"os"
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

		teams, err := cli.GetLoanTeams(v)
		for _, w := range teams {
			fmt.Println("team: ", w)
		}

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
	newLoans, err := cli.GetNewestLoans()
	fmt.Println("printing new loans.....")
	for _, v := range newLoans {
		fmt.Println(v)
	}
}
