package kiva

import (
    "io"
    "net/http"
    "net/url"
    "fmt"
    "encoding/json"
    "time"
    "path"
    "errors"
    "strconv"

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

type Lender struct {
  ID int `json: "lender_id"`
  Name string `json: "name"`
  Whereabouts string `json: "whereabots"`
  CountryCode string `json: "country_code"`
  Image
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


type doer interface {
  Do(*http.Request) (*http.Response, error)
}

type Client struct {
  baseURL url.URL
  doer doer
}

type Config struct {
  BaseURL string
}

func New(config *Config) *Client {
  if config == nil || config.BaseURL == "" {
    return &Client{
        baseURL: url.URL{
          Scheme: "http",
          Host: "api.kivaws.org",
        },
      doer: &http.Client{
        Timeout: time.Second * 10,
      },
    }
  }

  baseURL, err := url.Parse(config.BaseURL)
  if err != nil {
    panic(fmt.Sprintf("cannot parse base URL: %q (%v)", config.BaseURL, err))
  }
  return &Client{
    baseURL: *baseURL,
    doer: &http.Client{},
  }
}

// make the actual http request and return the http response
func (c *Client) raw(method string, urlpath string, query url.Values, body io.Reader) (*http.Response, error) {
  url := c.baseURL
  // FIXME: ".json" should probably be an option:
  urlpath += ".json"
  url.Path = path.Join(c.baseURL.Path, urlpath)
  fmt.Println(url.String())
  url.RawQuery = query.Encode()
  req, err := http.NewRequest(method, url.String(), body)
  if err != nil {
    return nil, err
  }
  return c.doer.Do(req)
}

// decode json from http response and return as interface understood
// by caller
func (c *Client) do(method string, urlpath string, query url.Values, body io.Reader, v interface{}) error {
  resp, err := c.raw(method, urlpath, query, body)
  if err != nil {
    return fmt.Errorf("error making request: %s", err)
  }
  decode := json.NewDecoder(resp.Body)
  if err = decode.Decode(&v); err !=nil {
    return fmt.Errorf("cannot decode your dumb response: %s", err)
  }
  return nil
}

type PagedLoansResponse struct {
  Paging PagingData `json: "paging"`
  Loans  []Loan     `json: "loans"`
}

type UnpagedLoansResponse struct {
  Loans []Loan `json: "loans"`
}

func (c *Client) GetLoansByID(ids ...int) ([]Loan, error) {
  // not sure whether requesting 50 loan IDs will return paged results

  var baseUrl = "/v1/loans/"
  var url string
  var loans []Loan
  if len(ids) == 0 {
    return loans, errors.New("No Loan IDs passed")
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
  
  var lr UnpagedLoansResponse
  err := c.do("GET", baseUrl + url, nil, nil, &lr)
  if err != nil {
    return nil, err
  }

  return lr.Loans, nil
}

func (c *Client) GetLoanJournalEntries(id int) {
  // not sure this is even implemented... they don't seem
  // to publish a schema for it either
  //baseUrl := fmt.Sprintf("/v1/loans/%d/journal_entries",id)
}

type PagedLendersResponse struct {
  Paging PagingData `json: "paging"`
  Lenders []Lender `json: "lenders"`
}

func (c *Client) GetLoanLenders(id int) ([]Lender, error) {

  baseUrl := fmt.Sprintf("/v1/loans/%d/lenders",id)
  var lr PagedLendersResponse

  err := c.do("GET", baseUrl, nil, nil, &lr)
  if err != nil {
    return nil, err
  }

  return lr.Lenders, nil
}