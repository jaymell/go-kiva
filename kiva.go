package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type Loan struct {
	name string
	location Location
	posted_date string
	activity string
	id int
	use string
	description Description
	funded_amount int
	partner_id int
	image Image
	borrower_count int
	loan_amount int
	status string
	sector string
}

type Image struct {
	template_id int
	id int
}

type Description struct {
	languages []string
}

type Location struct {
	country string
	geo string
	town string

}

func GetLoan(Url string) (Loan, error) {

	var l Loan
	resp, err := http.Get(Url + ".json")
	if err != nil {
		return l, err
	}
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return l, err
	}
	return l, nil
}

func GetRawJson(Url string) {

	var dat map[string]interface{}

	resp, err := http.Get(Url + ".json")
	if err != nil {
		fmt.Println(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bs, &dat)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range dat {
		fmt.Println(k)
		fmt.Println(v)
	}
}

func main() {
	urlBase := "http://api.kivaws.org/v1"
	url := "/loans/newest"

	//GetRawJson(urlBase + url)

	loan, err := GetLoan(urlBase + url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(loan.name)
}

