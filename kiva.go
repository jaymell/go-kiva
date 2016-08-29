package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Loan struct {
	Name                   string              `json:"name"`
	Location               LocationData        `json: "Location"`
	PostedDate             int                 `json: "posted_date"`
	Activity               string              `json: "activity"`
	Id                     int                 `json: "id"`
	Use                    string              `json: "use"`
	Desc                   Description         `json: "description"`
	FundedAmount           int                 `json: "funded_amount"`
	PartnerId              int                 `json: "partner_id"`
	Image                  Image               `json: "image"`
	BorrowerCount          int                 `json: "borrower_count"`
	LoanAmount             int                 `json: "loan_amount"`
	Status                 string              `json: "status"`
	Sector                 string              `json: "sector"`
	Expiration             int                 `json: "planned_expiration_date"`
	BonusCreditEligibility bool                `json: "bonus_credit_eligibility"`
	Tags                   []map[string]string `json: "tags"`
	BasketAmount           int                 `json: "basket_amount"`
}

type Image struct {
	TemplateId int `json: "template_id"`
	Id         int `json: "id"`
}

type Description struct {
	Languages []string `json: "languages"`
}

type LocationData struct {
	Country string            `json: "country"`
	Geo     map[string]string `json: "geo"`
	Town    string            `json: "town"`
}

type PagingData struct {
	Total    int `json: "total"`
	Page     int `json: "page"`
	PageSize int `json: "page_size"`
	Pages    int `json: "pages"`
}

type PagedLoansResponse struct {
	Paging PagingData `json: "paging"`
	Loans  []Loan     `json: "loans"`
}

type UnpagedLoansResponse struct {
	Loans []Loan `json: "loans"`
}

func Client() *client.Client {
	return client.New(&clientConfig)
}

func GetResponse(url string) (*http.Response, error) {
	urlBase := "http://api.kivaws.org/v1"
	r, err := http.Get(urlBase + url + ".json")
	if err != nil {
		return r, err
	}
	return r, nil
}

func GetLoansById(ids ...int) ([]Loan, error) {
	// not sure whether requesting 50 loan IDs will return paged results

	var baseUrl = "/loans/"
	var url string
	var err error
	var loans []Loan
	if len(ids) == 0 {
		return loans, errors.New("No Loan Ids passed")
	}
	for i, v := range ids {
		if i == 0 {
			char := strconv.Itoa(v)
			url += char
		} else {
			char := strconv.Itoa(v)
			url += "," + char
		}
	}
	r, err := GetResponse(baseUrl + url)
	if err != nil {
		return nil, err
	}
	var lr UnpagedLoansResponse
	err = json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		fmt.Println("error decoding json", err)
		return nil, err
	}
	return lr.Loans, nil
}

func PrintRawLoansJson() {

	var dat map[string]interface{}

	url := "/loans/newest"
	resp, err := GetResponse(url)
	if err != nil {
		fmt.Println(err)
		return
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
	//PrintRawLoansJson()

	loans, err := GetLoansById(1132720, 1128815)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("loans ", loans)
	for k, v := range loans {
		fmt.Println(k, v)
	}
}
