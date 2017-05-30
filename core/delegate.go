package core

import (
	"fmt"
	"net/http"
	"strconv"
)

//DelegateResponse data - received from api-call.
type DelegateResponse struct {
	Success        bool           `json:"success,omitempty"`
	Delegates      []DelegateData `json:"delegates,omitempty"`
	SingleDelegate DelegateData   `json:"delegate,omitempty"`
	TotalCount     int            `json:"totalCount,omitempty"`
}

//DelegateVoters struct to hold voters for a publicKey(delegate)
type DelegateVoters struct {
	Success  bool `json:"success"`
	Accounts []struct {
		Username  string `json:"username"`
		Address   string `json:"address"`
		PublicKey string `json:"publicKey"`
		Balance   string `json:"balance"`
	} `json:"accounts"`
}

//DelegateData holds parsed json from api calls. It is used in upper DelegateResponse struct
type DelegateData struct {
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

//DelegateQueryParams - when set, they are automatically added to get requests
type DelegateQueryParams struct {
	UserName  string `url:"username,omitempty"`
	PublicKey string `url:"publicKey,omitempty"`
	Offset    int    `url:"offset,omitempty"`
}

type DelegateDataProfit struct {
	Address         string
	VoteWeight      int
	VoteWeightShare float32
	EarnedAmmount   int
}

//Error interface function
func (e DelegateResponseError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v", e.Success, e.ErrorMessage)
}

//ListDelegates function returns list of delegtes. The top 51 delegates are returned
func (s *ArkClient) ListDelegates(params DelegateQueryParams) (DelegateResponse, *http.Response, error) {
	respData := new(DelegateResponse)
	respError := new(DelegateResponseError)
	resp, err := s.sling.New().Get("api/delegates").QueryStruct(&params).Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}

//GetDelegate function returns a delegate
func (s *ArkClient) GetDelegate(params DelegateQueryParams) (DelegateResponse, *http.Response, error) {
	respData := new(DelegateResponse)
	respError := new(DelegateResponseError)
	resp, err := s.sling.New().Get("api/delegates/get").QueryStruct(&params).Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}

//GetDelegateVoters function returns a delegate
func (s *ArkClient) GetDelegateVoters(params DelegateQueryParams) (DelegateVoters, *http.Response, error) {
	respData := new(DelegateVoters)
	respError := new(DelegateResponseError)
	resp, err := s.sling.New().Get("api/delegates/voters").QueryStruct(&params).Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}

//GetDelegateVoteWeight function returns a summary of ARK voted for selected delegate
func (s *ArkClient) GetDelegateVoteWeight(params DelegateQueryParams) (int, *http.Response, error) {
	respData := new(DelegateVoters)
	respError := new(DelegateResponseError)
	resp, err := s.sling.New().Get("api/delegates/voters").QueryStruct(&params).Receive(respData, respError)
	if err == nil {
		err = respError
	}

	//calculating vote weight
	balance := 0
	if respData.Success {
		for _, element := range respData.Accounts {
			intBalance, _ := strconv.Atoi(element.Balance)
			balance += intBalance
		}
	}

	return balance, resp, err
}

func (s *ArkClient) CalculateVotersProfit(params DelegateQueryParams) ([]DelegateDataProfit, *http.Response, error) {
	respData := new(DelegateVoters)
	respError := new(DelegateResponseError)
	resp, err := s.sling.New().Get("api/delegates/voters").QueryStruct(&params).Receive(respData, respError)
	if err == nil {
		err = respError
	}

	//calculating vote weight
	votersProfit := []DelegateDataProfit{}
	balance := 0
	if respData.Success {
		//computing summ of all votes
		for _, element := range respData.Accounts {
			intBalance, _ := strconv.Atoi(element.Balance)
			balance += intBalance
		}

		//calculating
		for _, element := range respData.Accounts {
			deleProfit := DelegateDataProfit{
				Address: element.Address,
			}

			currentBalance, _ := strconv.Atoi(element.Balance)
			deleProfit.VoteWeight = currentBalance
			deleProfit.VoteWeightShare = float32(currentBalance / balance)

			votersProfit = append(votersProfit, deleProfit)
		}

	}

	return votersProfit, resp, err
}
