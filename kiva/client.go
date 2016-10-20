package kiva

/*
paged methods:
/lenders/:lender_id/teams
/lenders/newest
/lenders/search
/loans/:id/journal_entries
/loans/:id/lenders
/loans/:id/teams
/loans/newest
/loans/search
/my/loans
/my/loans/:ids
/my/teams
/partners
/teams/:id/lenders
/teams/:id/loans
/teams/search
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type Loan struct {
	Name                   string              `json:"name"`
	Location               LocationData        `json: "Location"`
	PostedDate             int                 `json: "posted_date"`
	Activity               string              `json: "activity"`
	ID                     int                 `json: "id"`
	Use                    string              `json: "use"`
	Desc                   Description         `json: "description"`
	FundedAmount           int                 `json: "funded_amount"`
	PartnerID              int                 `json: "partner_id"`
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

type Team struct {
	Category       string `json: "category"`
	Descrption     string `json: "description"`
	ID             int    `json: "id"`
	Image          Image  `json: "image"`
	LoanBecause    string `json: "loan_because"`
	LoanCount      int    `json: "loan_count"`
	LoanedAmount   int    `json: "loaned_amount"`
	MemberCount    int    `json: "member_count"`
	MembershipType int    `json: "membership_type"`
	Name           string `json: "name"`
	ShortName      string `json: "short_name"`
	// FIXME: string for now -- example: 2013-11-03T13:27:16Z
	TeamSince string `json: "team_since"`
	// FIXME: this could probably be url type:
	WebsiteURL  string `json: "website_url"`
	Whereabouts string `json: "whereabouts"`
}

type Lender struct {
	ID          int    `json: "lender_id"`
	Name        string `json: "name"`
	Whereabouts string `json: "whereabots"`
	CountryCode string `json: "country_code"`
	Image
}

type Repayment struct {
	// can't find data on this
}

type Image struct {
	TemplateID int `json: "template_id"`
	ID         int `json: "id"`
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

type PagedLoanResponse struct {
	Paging PagingData `json: "paging"`
	Items  []Loan     `json: "loans"`
}

type PagedLenderResponse struct {
	Paging  PagingData `json: "paging"`
	Items []Lender   `json: "lenders"`
}

// type PagedLoanRepaymentsResponse struct {
//   Paging PagingData `json: "paging"`

// }

type PagedLoanTeamResponse struct {
	Paging PagingData `json: "paging"`
	Items []Team `json: "teams"`
}

type UnpagedLoanResponse struct {
	Loans []Loan `json: "loans"`
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	baseURL *url.URL
	doer    doer
	appID string
}

type Config struct {
	BaseURL string
	AppID string
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL: "http://api.kivaws.org",
		AppID: "",		
	}
}

func New(config *Config) *Client {
	defaultConfig := DefaultConfig()
	if config == nil {
		config = defaultConfig
	} 

    if config.BaseURL == "" {
		config.BaseURL = defaultConfig.BaseURL
    }
    
    baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		panic(fmt.Sprintf("cannot parse base URL: %q (%v)", config.BaseURL, err))
	}

	return &Client{
		baseURL: baseURL,
		doer:    &http.Client{
				Timeout: time.Second * 10,
			},
		appID: config.AppID,
	}
}

// make the actual http request and return the http response
func (c *Client) raw(method string, urlpath string, query url.Values, body io.Reader) (*http.Response, error) {
	url := *c.baseURL
	// FIXME: ".json" should probably be an option:
	urlpath += ".json"
	url.Path = path.Join(url.Path, urlpath)
	url.RawQuery = query.Encode()
	fmt.Println(url.String())
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	return c.doer.Do(req)
}

// decode json from http response and return as interface understood
// by caller
func (c *Client) do(method string, urlpath string, query url.Values, body io.Reader, v interface{}) error {
	// add application id if defined:
	newQuery := url.Values{}
	if c.appID != "" {
		if query == nil {
			query = newQuery
		}
		query.Set("app_id", c.appID)		
	}

	resp, err := c.raw(method, urlpath, query, body)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}
	decode := json.NewDecoder(resp.Body)
	if err = decode.Decode(&v); err != nil {
		return fmt.Errorf("cannot decode json: %s", err)
	}
	return nil
}



// wraps "do" to handle paged requests
func (c *Client) doPaged(method string, urlpath string, query url.Values, body io.Reader, v interface{}, numPages int) (error) {

	var pr interface{}

	if query == nil {
		query = url.Values{}
	}

	if numPages < 0 {
		return fmt.Errorf("less than zero is unacceptable") 
	}

  	// get the first page
	err := c.do(method, urlpath, query, body, pr)
	if err != nil {
		return err
	}

	var items interface{}
	var newpr interface{}

	switch pr := pr.(type) {
		default:
			fmt.Errorf("Unexpected type %T", pr)
		case PagedLoanResponse:
			newpr = pr
			items = make([]Loan, 0, pr.Paging.Total)
		case PagedLenderResponse:
			newpr = pr
			items = make([]Lender, 0, pr.Paging.Total)
		case PagedLoanTeamResponse:
			newpr = pr
			items = make([]Team, 0, pr.Paging.Total)
	}

    if newpr.Paging.Pages == 1 {
    	return pr.Items, nil
    } 
    
    // get all pages if zero -- maybe a bad idea:
    if numPages == 0 {
    	numPages = newpr.Paging.Pages
    }

	items = append(items, newpr.Items...)
	
	for i:=2; i <= numPages; i++ {
     	query.Set("page", strconv.Itoa(i)) 
		err := c.do("GET", urlpath, query, nil, &newpr)
		if err != nil {
			return nil, err
		}
		items = append(items, pr.Items...)
	}

	return items, nil
}

func (c *Client) GetLoansByID(loanIDs ...int) ([]Loan, error) {
	// not sure whether requesting 50 loan IDs will return paged results

	var baseURL = "/v1/loans"
	var url string
	var loans []Loan
	if len(loanIDs) == 0 {
		return loans, errors.New("No Loan IDs passed")
	}
	for i, v := range loanIDs {
		if i == 0 {
			char := strconv.Itoa(v)
			url += char
		} else {
			char := strconv.Itoa(v)
			url += "," + char
		}
	}

	var lr UnpagedLoanResponse
	err := c.do("GET", baseURL+url, nil, nil, &lr)
	if err != nil {
		return nil, err
	}

	return lr.Loans, nil
}

func (c *Client) GetLoanJournalEntries(loanID int) {
	// not sure this is even implemented... they don't seem
	// to publish a schema for it either
	//baseURL := fmt.Sprintf("/v1/loans/%d/journal_entries", loanID)
}



// FIXME: need to return ALL pages
func (c *Client) GetLoanLenders(loanID int) ([]Lender, error) {

	baseURL := fmt.Sprintf("/v1/loans/%d/lenders", loanID)
	var lr PagedLenderResponse

	err := c.do("GET", baseURL, nil, nil, &lr)
	if err != nil {
		return nil, err
	}
	// FIXME: need to return ALL pages
	return lr.Lenders, nil
}



// func (c *Client) GetLoanRepayments(loanID int) {
//  // can't find data on this
// }

func (c *Client) GetSimilarLoans(loanID int) ([]Loan, error) {
	baseURL := fmt.Sprintf("/v1/loans/%d/similar", loanID)
	var lr UnpagedLoanResponse

	err := c.do("GET", baseURL, nil, nil, &lr)
	if err != nil {
		return nil, err
	}

	return lr.Loans, nil
}




// FIXME: a lot of common code 
// for paged responses can be factored
// FIXME: be able to pass options
// FIXME: don't just return all pages at once
func (c *Client) GetLoanTeams(loanID int) ([]Team, error) {
	baseURL := fmt.Sprintf("/v1/loans/%d/teams", loanID)
	var pr PagedLoanTeamsResponse

	numPages := 1 // set initial value of number of pages to iterate through
	err := c.do("GET", baseURL, nil, nil, &pr)
	if err != nil {
		return nil, err
	}

    if pr.Paging.Pages < 2 {
    	return pr.Teams, nil
    }
    numPages = pr.Paging.Pages

	teams := make([]Team, pr.Paging.Total)
	copy(teams, pr.Teams)
	query := url.Values{}

	for i:=2; i <= numPages; i++ {
     	query.Set("page", strconv.Itoa(i)) 
		err := c.do("GET", baseURL, query, nil, &pr)
		if err != nil {
			return nil, err
		}
		numPages = pr.Paging.Pages // update numPages based on subsequent responses
		copy(teams, pr.Teams)
	}

	return teams, nil
}

// FIXME: be able to pass options
// FIXME: need to be able to handle possible changes in total
	   // number of records when getting a large paged response ??
func (c *Client) GetNewestLoans(numPages int) ([]Loan, error) {
  baseURL := "/v1/loans/newest"
  return c.doPagedLoan("GET", baseURL, nil, nil, numPages)
}

