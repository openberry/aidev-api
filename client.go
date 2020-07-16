package aidev

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// API endpoints
const (
	testingURL = "http://aidev.work/truth/Connect/v1"
)

// Available configuration options, if not provided sane values will be
// used by default
type Options struct {
	// Username. Used to automatically renew tokens when required.
	User string

	// Used to automatically renew tokens when required.
	Password string

	// JWT API access credentials.
	AccessToken string

	// Time to wait for requests, in seconds
	Timeout uint

	// Time to maintain open the connection with the service, in seconds
	KeepAlive uint

	// Maximum network connections to keep open with the service
	MaxConnections uint

	// User agent value to report to the service
	UserAgent string

	// Produce trace output of requests and responses.
	Debug bool
}

// Main service handler
type Client struct {
	BaseAPI *baseAPI
	token     string
	options   Options
	conn      *http.Client
}

// Network request options
type requestOptions struct {
	method   string
	endpoint string
	data     map[string]interface{}
}

// Return sane default configuration values
func defaultOptions() *Options {
	return &Options{
		Timeout:        30,
		KeepAlive:      600,
		MaxConnections: 100,
		UserAgent:      "aidev-lib/0.1.0",
	}
}

// NewClient will construct a usable service handler using the provided API key and
// configuration options, if 'nil' options are provided default sane values will
// be used
func NewClient(options *Options) *Client {
	// If no options are provided, use default sane values
	if options == nil {
		options = defaultOptions()
	}

	// Configure base HTTP transport
	t := &http.Transport{
		MaxIdleConns:        int(options.MaxConnections),
		MaxIdleConnsPerHost: int(options.MaxConnections),
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(options.Timeout) * time.Second,
			KeepAlive: time.Duration(options.KeepAlive) * time.Second,
		}).DialContext,
	}

	// Setup main client
	client := &Client{
		options: *options,
		conn: &http.Client{
			Transport: t,
			Timeout:   time.Duration(options.Timeout) * time.Second,
		},
	}

	// Load API modules
	client.BaseAPI = &baseAPI{cl: client}
	return client
}

// SetToken adjust the API access token used by the client instance.
func (i *Client) SetToken(token string) {
	i.token = token
}

// RenewToken request a new access token and set it on the client instance.
func (i *Client) RenewToken() error {
	response, err := i.request(&requestOptions{
		method:   "POST",
		endpoint: "/getToken",
		data: map[string]interface{}{
			"nick": i.options.User,
			"psw":  i.options.Password,
		},
	})
	if err != nil {
		return err
	}
	token, ok := response["token"]
	if !ok {
		return errors.New(errInvalidResponse)
	}
	i.token = token.(string)
	return nil
}

// Dispatch a network request to the service
func (i *Client) request(r *requestOptions) (map[string]interface{}, error) {
	// Build HTTP request
	form := url.Values{}
	for k, v := range r.data {
		form.Set(k, fmt.Sprintf("%s", v))
	}
	req, _ := http.NewRequest(r.method, testingURL + r.endpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Add user-agent header
	if i.options.UserAgent != "" {
		req.Header.Add("User-Agent", i.options.UserAgent)
	}

	// Dump request
	if i.options.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			fmt.Println("=== request ===")
			fmt.Printf("%s\n", dump)
		}
	}

	// Execute request
	res, err := i.conn.Do(req)
	if res != nil {
		// Properly discard request content to be able to reuse the connection
		defer func() {
			_, _ = io.Copy(ioutil.Discard, res.Body)
			_ = res.Body.Close()
		}()
	}

	// Dump response
	if i.options.Debug {
		dump, err := httputil.DumpResponse(res, true)
		if err == nil {
			fmt.Println("=== response ===")
			fmt.Printf("%s\n", dump)
		}
	}

	// Network level errors
	if err != nil {
		return nil, err
	}

	// Get response contents
	body, err := ioutil.ReadAll(res.Body)

	// Application level errors
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unsucessfull request: %s", res.Status))
	}

	// Decode response
	wrapper := map[string]interface{}{}
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return nil, errors.New(fmt.Sprintf("non JSON content returned: %s", body))
	}
	if err, ok := wrapper["error"]; ok {
		if err.(string) != "" {
			return nil, errors.New(err.(string))
		}
	}
	return wrapper, nil
}
