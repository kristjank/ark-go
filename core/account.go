package core

import (
	"fmt"
	"net/http"
)

//DelegateResponse data - received from api-call.
type DelegateResponse struct {
	Success    bool           `json:"success"`
	Delegates  []DelegateData `json:"delegates"`
	TotalCount int            `json:"totalCount"`
}

//DelegateData holds parsed json from api calls. It is used in upper DelegateResponse struct
type DelegateData []struct {
	Username       string  `json:"username"`
	Address        string  `json:"address"`
	PublicKey      string  `json:"publicKey"`
	Vote           string  `json:"vote"`
	Producedblocks int     `json:"producedblocks"`
	Missedblocks   int     `json:"missedblocks"`
	Rate           int     `json:"rate"`
	Approval       float64 `json:"approval"`
	Productivity   float64 `json:"productivity"`
}

//DelegateResponseError struct to hold error response
type DelegateResponseError struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}

//Error interface function
func (e DelegateResponseError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v", e.Success, e.ErrorMessage)
}

//ListDelegates function returns list of peers from ArkNode
func (s *ArkClient) ListDelegates() (DelegateResponse, *http.Response, error) {
	respData := new(DelegateResponse)
	respError := new(DelegateResponseError)
	resp, err := s.sling.New().Get("api/delegates").Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}
