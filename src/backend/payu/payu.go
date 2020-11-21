package payugo

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type MerchantConfig struct {
	ClientID     string
	ClientSecret string
	PosID        string
}

type Client struct {
	client      *http.Client
	baseURL     *url.URL
	cfg         MerchantConfig
	accessToken string
}

func NewClient(client *http.Client, baseURL string, cfg MerchantConfig) (*Client, error) {
	if client == nil {
		client = http.DefaultClient
	}

	// ensure the baseURL contains a trailing slash so that all paths are preserved in later calls
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &Client{
		client:  client,
		baseURL: parsedBaseURL,
		cfg:     cfg,
	}, nil
}

type Language string

const (
	LanguagePL Language = "pl" // polish
	LanguageEN Language = "en" // english
	LanguageCS Language = "cs" // czech
	LanguageDE Language = "de" // deutsch
)

type CurrencyCode string

const (
	CurrencyCodePLN = "PLN"
	CurrencyCodeEUR = "EUR"
	CurrencyCodeUSD = "USD"
)

type Order struct {
	NotifyURL     string       `json:"notifyUrl"`
	CustomerIP    string       `json:"customerIp"`
	MerchantPosID string       `json:"merchantPosId"` // value will be overwritten with given value from config
	Description   string       `json:"description"`
	CurrencyCode  CurrencyCode `json:"currencyCode"`
	TotalAmount   string       `json:"totalAmount"`
	ContinueURL   string       `json:"continueUrl"`
	ExtOrderID    string       `json:"extOrderId"` // user will be redirect here after payment
	Buyer         Buyer        `json:"buyer"`
	Products      []Product    `json:"products"`
	OrderInfo
}

// Validate checks if all necessary fields were provided, should be called right before sending http request
func (o Order) validate() error {
	if o.CustomerIP == "" {
		return fmt.Errorf("%w, no CustomerIP", ErrOrderInvalid)
	} else if o.MerchantPosID == "" {
		return fmt.Errorf("%w, no MerchantPosID", ErrOrderInvalid)
	} else if o.Description == "" {
		return fmt.Errorf("%w, no Description", ErrOrderInvalid)
	} else if o.CurrencyCode == "" {
		return fmt.Errorf("%w, no CurrencyCode", ErrOrderInvalid)
	} else if o.TotalAmount == "" {
		return fmt.Errorf("%w, no TotalAmount", ErrOrderInvalid)
	} else if o.Products == nil || len(o.Products) == 0 {
		return fmt.Errorf("%w, no Products", ErrOrderInvalid)
	}
	return nil
}

type Buyer struct {
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Language  Language `json:"language"`
}

type Product struct {
	Name      string `json:"name"`
	UnitPrice string `json:"unitPrice"`
	Quantity  string `json:"quantity"`
}

type OrderInfo struct {
	OrderID    string `json:"orderId"`
	ExtOrderID string `json:"extOrderId"`
}

type CreateOrderResponse struct {
	Status      StatusHolder `json:"status"`
	RedirectURI string       `json:"redirectUri"`
	OrderInfo
}

type RetrieveOrderResponse struct {
	Orders []Order      `json:"orders"`
	Status StatusHolder `json:"status"`
}

type StatusHolder struct {
	StatusCode Status `json:"statusCode"`
	StatusDesc string `json:"statusDesc"`
}

type StatusError struct {
	StatusHolder
	Code        string `json:"code"`
	CodeLiteral string `json:"codeLiteral"`
}

type ErrResponse struct {
	Err StatusError `json:"status"`
}

// we preserve that baseURL contains a trailing slash so endpoint shouldn't start with one
const (
	orderCreateRequestEndpoint = "api/v2_1/orders"
	orderRetrieveRequest       = "api/v2_1/orders/%s"
	authorizeEndpoint          = "pl/standard/user/oauth/authorize"
)

var (
	ErrInvalidStatus = errors.New("invalid status")
	ErrOrderInvalid  = errors.New("order invalid")
	ErrNoAccessToken = errors.New("you need to authorize client before using it")
)
