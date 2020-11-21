package payugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c Client) OrderCreateRequest(o *Order) (*CreateOrderResponse, error) {
	rel, err := url.Parse(orderCreateRequestEndpoint)
	if err != nil {
		return nil, err
	}
	o.MerchantPosID = c.cfg.PosID // set MerchantPosID with given one
	u := c.baseURL.ResolveReference(rel)

	//if err := o.validate(); err != nil {
	//	return nil, err
	//}
	body, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	if err := c.authorizeRequest(req); err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusFound {
		var errResp ErrResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w, with errorStatusCode: %s, codeLiteral: %s, statusDesc: %s",
			ErrInvalidStatus, errResp.Err.StatusCode, errResp.Err.CodeLiteral, errResp.Err.StatusDesc)

	}

	var orderResponse CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		return nil, err
	}
	return &orderResponse, nil
}

func (c Client) OrderRetrieveRequest(orderID string) (*RetrieveOrderResponse, error) {
	rel, err := url.Parse(fmt.Sprintf(orderRetrieveRequest, orderID))
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	if err := c.authorizeRequest(req); err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var errResp ErrResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w, with errorStatusCode: %s, codeLiteral: %s, statusDesc: %s",
			ErrInvalidStatus, errResp.Err.StatusCode, errResp.Err.CodeLiteral, errResp.Err.StatusDesc)
	}

	var orderResponse RetrieveOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		return nil, err
	}
	return &orderResponse, nil
}
