package payugo

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const _baseURL = "https://secure.snd.payu.com/"

func TestClient_Authorize(t *testing.T) {
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
		t.Error(err)
	}

	mockResp := AuthResponse{
		AccessToken: "8f79c971-195e-43f5-bd83-ad2104414acc",
		TokenType:   "bearer",
		ExpiresIn:   1000,
		GrantType:   "client_credentials",
	}
	b, err := json.Marshal(mockResp)
	if err != nil {
		t.Error(err)
	}
	httpmock.Activate()
	defer httpmock.Deactivate()
	httpmock.RegisterResponder(http.MethodPost, _baseURL+authorizeEndpoint,
		httpmock.NewBytesResponder(http.StatusOK, b))

	if err := c.Authorize(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, mockResp.AccessToken, c.accessToken)
}
