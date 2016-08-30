package kiva

import (
    "net/http"
    "net/url"
    "fmt"
    "encoding/json"
    "path"
    "time"
)

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
func (c *Client) raw(method string, path string, query url.Values, body io.Reader) (*http.Response, error) {
  url := c.BaseURL
  url.Path = path.Join(c.baseURL,path)
  url.RawQuery = query.Encode()
  req, err := http.NewRequest(method, url.String(), body)
  if err != nil {
    return nil, err
  }
  return c.doer.Do(req)
}

// decode json from http response and return as interface understood
// by caller
func (c *Client) do(method string, path string, query url.Values, body io.Reader, v interface{}) error {
  resp, err := c.raw(method, path, query, body)
  if err != nil {
    return fmt.Errorf("error making request: %s", err)
  }
  decode := json.NewDecoder(resp.Body)
  if err = decode.Decode(&v); err !=nil {
    return fmt.Errorf("cannot decode your dumb response: %s", err)
  }
  return nil
}

func (c *Client) GetLoansById(ids ...int) ([]Loan, error) {
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
  
  var lr UnpagedLoansResponse
  r, err := c.do("GET", baseUrl + url, nil, nil, &lr)
  if err != nil {
    return nil, err
  }
  
  return lr.Loans, nil
}

