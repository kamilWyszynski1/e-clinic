package payugo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/magiconair/properties/assert"
)

const _extOrderID = "externalOrderID"

func TestClient_OrderCreateRequest(t *testing.T) {
	c, err := NewClient(
		http.DefaultClient,
		_baseURL,
		MerchantConfig{
			ClientID:     "123",
			ClientSecret: "123",
			PosID:        "123",
		},
	)
	if err != nil {
		t.Error(err)
	}
	c.accessToken = "token"

	mockResp := CreateOrderResponse{
		Status:      StatusHolder{StatusCode: StatusSuccess},
		RedirectURI: "redirURI",
		OrderInfo: OrderInfo{
			OrderID:    "orderID",
			ExtOrderID: _extOrderID},
	}
	b, err := json.Marshal(mockResp)
	if err != nil {
		t.Error(err)
	}
	httpmock.Activate()
	defer httpmock.Deactivate()
	httpmock.RegisterResponder(http.MethodPost, _baseURL+orderCreateRequestEndpoint,
		httpmock.NewBytesResponder(http.StatusOK, b))

	resp, err := c.OrderCreateRequest(&Order{
		NotifyURL:     "notifURL",
		CustomerIP:    "",
		MerchantPosID: "",
		Description:   "test order",
		CurrencyCode:  CurrencyCodeEUR,
		TotalAmount:   "1500",
		ExtOrderID:    _extOrderID,
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.ExtOrderID, _extOrderID)
}

func TestClient_OrderCreateRequestIntegration(t *testing.T) {
	//t.Skip()
	c, err := NewClient(
		http.DefaultClient,
		_baseURL,
		MerchantConfig{
			ClientID:     "398268",
			ClientSecret: "880487191465ca9418fafcd9c0a019e6",
			PosID:        "398268",
		},
	)
	if err != nil {
		t.Error(err)
	}
	if err := c.Authorize(); err != nil {
		t.Fatal(err)
	}

	resp, err := c.OrderCreateRequest(&Order{
		NotifyURL:    "https://webhook.site/9534fb13-a140-499f-9a2b-faaf61c6c400",
		CustomerIP:   "127.0.0.1",
		Description:  "test order",
		CurrencyCode: CurrencyCodePLN,
		TotalAmount:  "1",
		Buyer: Buyer{
			Email:     "kamil.wyszynski.97@gmail.com",
			Phone:     "",
			FirstName: "Kamil",
			LastName:  "Wyszyński",
			Language:  LanguagePL,
		},
		Products: []Product{
			{
				Name:      "Test product",
				UnitPrice: "1",
				Quantity:  "1",
			},
		},
	})
	fmt.Println(err)
	fmt.Println(resp)
	fmt.Println(resp.OrderID)
}

func TestClient_OrderRetrieveRequest(t *testing.T) {
	//t.Skip()
	c, err := NewClient(
		http.DefaultClient,
		_baseURL,
		MerchantConfig{
			ClientID:     "398268",
			ClientSecret: "880487191465ca9418fafcd9c0a019e6",
			PosID:        "398268",
		},
	)
	if err != nil {
		t.Error(err)
	}
	if err := c.Authorize(); err != nil {
		t.Fatal(err)
	}
	resp, err := c.OrderRetrieveRequest("HDHRSQ6Q3N201121GUEST000P01")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}
