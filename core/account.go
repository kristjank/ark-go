package core

import "net/http"

//AccountResponse structure
type AccountResponse struct {
	Success bool        `json:"success"`
	Account AccountData `json:"account"`
}

//AccountData structure
type AccountData struct {
	Address              string        `json:"address"`
	UnconfirmedBalance   string        `json:"unconfirmedBalance"`
	Balance              string        `json:"balance"`
	PublicKey            string        `json:"publicKey"`
	UnconfirmedSignature int           `json:"unconfirmedSignature"`
	SecondSignature      int           `json:"secondSignature"`
	SecondPublicKey      interface{}   `json:"secondPublicKey"`
	Multisignatures      []interface{} `json:"multisignatures"`
	UMultisignatures     []interface{} `json:"u_multisignatures"`
}

//AccountQueryParams structure
type AccountQueryParams struct {
	Address string `url:"address,omitempty"`
}

//GetAccount by address function returns list of peers from ArkNode
func (s *ArkClient) GetAccount(params AccountQueryParams) (AccountResponse, *http.Response, error) {
	accResponse := new(AccountResponse)
	accResponseError := new(ArkApiResponseError)

	resp, err := s.sling.New().Get("api/accounts").QueryStruct(&params).Receive(accResponse, accResponseError)
	if err == nil {
		err = accResponseError
	}

	return *accResponse, resp, err
}
