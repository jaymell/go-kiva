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
    "log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/jpillora/backoff"
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

type pagingData struct {
	Total    int `json: "total"`
	Page     int `json: "page"`
	PageSize int `json:"page_size"`
	Pages    int `json: "pages"`
}

type Pageable interface {
	Paging() pagingData
}

type Pager struct {
	PagingData pagingData `json:"paging"`
}

func (p Pager) Paging() pagingData {
	return p.PagingData
}

type PagedLoanResponse struct {
	Pager
	Loans []Loan `json: "loans"`
}

type PagedLenderResponse struct {
	Pager
	Lenders []Lender `json: "lenders"`
}

// type PagedLoanRepaymentsResponse struct {
//   Paging pagingData `json: "paging"`

// }

type PagedTeamResponse struct {
	Pager
	Teams []Team `json: "teams"`
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
	appID   string
	backoff *backoff.Backoff	
}

type Config struct {
	BaseURL string
	AppID   string
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL: "http://api.kivaws.org",
		AppID:   "",
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

	b := &backoff.Backoff{
	    Jitter: true,
	    Min: 1 * time.Second,
	    Max: 1 * time.Minute,
	}

	return &Client{
		baseURL: baseURL,
		doer: &http.Client{
			Timeout: time.Second * 10,
		},
		appID: config.AppID,
		backoff: b,
	}
}

// make the actual http request and return the http response
func (c *Client) raw(method string, urlpath string, query url.Values, body io.Reader) (*http.Response, error) {
	url := *c.baseURL
	// FIXME: ".json" should probably be an option:
	urlpath += ".json"
	url.Path = path.Join(url.Path, urlpath)
	url.RawQuery = query.Encode()

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	// this will retry forever if hit rate limit; 
	// for now, leaving it up to caller to
	// decide when to give up on a request:	
	var resp *http.Response
	for {
		resp, err = c.doer.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.Status == "403 Forbidden (Rate Limit Exceeded)" {
			d := c.backoff.Duration()
			log.Println("Rate throttled. Backing off: ", d)
			time.Sleep(d)
			continue
		}
		c.backoff.Reset()
		break
	}

	//fmt.Println("X-RateLimit-Overall-Limit: ", resp.Header.Get("X-RateLimit-Overall-Limit"))
	//fmt.Println("X-RateLimit-Overall-Remaining: ", resp.Header.Get("X-RateLimit-Overall-Remaining"))
	//fmt.Println("X-RateLimit-Specific-Limit: ", resp.Header.Get("X-RateLimit-Specific-Limit"))
	//fmt.Println("X-RateLimit-Specific-Remaining: ", resp.Header.Get("X-RateLimit-Specific-Remaining"))

	return resp, err	
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

type ChannelResult struct {
    Val *http.Response
    Err error
}

// thread-safe "do" to handle paged requests
func (c *Client) doPaged(urlpath string, query url.Values, pr Pageable, numPages int) ([]Pageable, error) {
	// FIXME: reimplements c.do, gets 1st 
	// page of response and decodes it, uses that info to determine
	// how many subsequent pages there will be, then gets all of those
	// raw responses in a go routine, then decodes all the raw
	// responses... surely this be done more succinctly

	if numPages < 0 {
		return nil, fmt.Errorf("less than zero is unacceptable")
	}
	
	newQuery := url.Values{}
	if c.appID != "" {
		if query == nil {
			query = newQuery
		}
		query.Set("app_id", c.appID)
	}

	raw, err := c.raw("GET", urlpath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err)
	}

	// have to decode first response so we know how many pages:
	decode := json.NewDecoder(raw.Body)
	if err = decode.Decode(&pr); err != nil {
		return nil, fmt.Errorf("cannot decode json: %s", err)
	}

	resp := make([]Pageable, 1)
	resp[0] = pr
	paging := resp[0].Paging()

	if numPages == 1 || paging.Pages <= 1 {
		return resp, nil
	}

	// get all pages if zero:
	if numPages == 0 {
		numPages = paging.Pages
	}

	// array of raw responses which will
	// be decoded AFTER all have been collected:

	ch := make(chan ChannelResult)
	for i := 2; i <= numPages; i++ {
		go func(i int) {
			q := url.Values{}
			for k,v := range query {
  				q[k] = v
			}
			q.Set("page", strconv.Itoa(i))
			raw, err := c.raw("GET", urlpath, query, nil)
			if err != nil {
				ch <- ChannelResult{Val: nil, Err: err }
			}
			ch <- ChannelResult{Val: raw, Err: nil }
		}(i)
	}
	
	var rawArr []*http.Response
	for {
		select {
			case r := <- ch:
				if r.Err != nil {
					return nil, err
				}
				rawArr = append(rawArr, r.Val)

		}
		if len(rawArr) == numPages-1 { 
			break
		}
	}

	respArr := make([]Pageable, 1)
	copy(respArr, resp)
	for _, val := range rawArr {
		decode := json.NewDecoder(val.Body)
		if err = decode.Decode(&pr); err != nil {
			return nil, fmt.Errorf("cannot decode json: %s", err)
		}
		respArr = append(respArr, pr)	
	}

	return respArr, nil
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

func (c *Client) GetSimilarLoans(loanID int) ([]Loan, error) {

	// FIXME -- allow to specify count

	baseURL := fmt.Sprintf("/v1/loans/%d/similar", loanID)
	var lr UnpagedLoanResponse

	err := c.do("GET", baseURL, nil, nil, &lr)
	if err != nil {
		return nil, err
	}

	return lr.Loans, nil
}

func (c *Client) GetPagedLoanResponse(urlpath string, query url.Values, numPages int) ([]Loan, error) {
	var respType PagedLoanResponse
	prArray, err := c.doPaged(urlpath, query, &respType, numPages)
	if err != nil {
		return nil, err
	}
	paging := prArray[0].Paging()
	loans := make([]Loan, len(prArray)*paging.PageSize)
	iter := 0
	for _, v := range prArray {
		pr := v.(*PagedLoanResponse)
		for _, loan := range pr.Loans {
			loans[iter] = loan
			iter++
		}
	}
	return loans, nil
}

func (c *Client) GetNewestLoans(numPages int) ([]Loan, error) {
	baseURL := "/v1/loans/newest"
	loans, err := c.GetPagedLoanResponse(baseURL, nil, numPages)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (c *Client) GetPagedLenderResponse(urlpath string, query url.Values, numPages int) ([]Lender, error) {
	var respType PagedLenderResponse
	prArray, err := c.doPaged(urlpath, query, &respType, numPages)
	if err != nil {
		return nil, err
	}
	paging := prArray[0].Paging()
	lenders := make([]Lender, len(prArray)*paging.PageSize)
	iter := 0
	for _, v := range prArray {
		pr := v.(*PagedLenderResponse)
		for _, lender := range pr.Lenders {
			lenders[iter] = lender
			iter++
		}
	}
	return lenders, nil
}

func (c *Client) GetLoanLenders(loanID int) ([]Lender, error) {

	baseURL := fmt.Sprintf("/v1/loans/%d/lenders", loanID)
	lenders, err := c.GetPagedLenderResponse(baseURL, nil, 0)
	if err != nil {
		return nil, err
	}
	return lenders, nil
}

func (c *Client) GetPagedTeamResponse(urlpath string, query url.Values, numPages int) ([]Team, error) {
	var respType PagedTeamResponse
	prArray, err := c.doPaged(urlpath, query, &respType, numPages)
	if err != nil {
		return nil, err
	}
	paging := prArray[0].Paging()
	teams := make([]Team, len(prArray)*paging.PageSize)
	iter := 0
	for _, v := range prArray {
		pr := v.(*PagedTeamResponse)
		for _, team := range pr.Teams {
			teams[iter] = team
			iter++
		}
	}
	return teams, nil
}

func (c *Client) GetLoanTeams(loanID int) ([]Team, error) {
	baseURL := fmt.Sprintf("/v1/loans/%d/teams", loanID)
	teams, err := c.GetPagedTeamResponse(baseURL, nil, 0)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// func (c *Client) GetLoanRepayments(loanID int) {
//  // can't find data on this
// }

// func (c *Client) GetLoanJournalEntries(loanID int) {
// 	// not sure this is even implemented... they don't seem
// 	// to publish a schema for it either
// 	//baseURL := fmt.Sprintf("/v1/loans/%d/journal_entries", loanID)
// }
