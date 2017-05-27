package core

import "fmt"

const (
	MAINNET = iota
	DEVNET
)

type ArkEnvParams struct {
	Success bool    `json:"success"`
	Network Network `json:"network"`
	Fees    Fees    `json:"fees"`
}

//Fees constant parameters for various ArkCoin configurations
type Fees struct {
	Send            int64 `json:"send"`
	Vote            int64 `json:"vote"`
	Secondsignature int64 `json:"secondsignature"`
	Delegate        int64 `json:"delegate"`
	Multisignature  int64 `json:"multisignature"`
}

type Network struct {
	Nethash  string `json:"nethash"`
	Token    string `json:"token"`
	Symbol   string `json:"symbol"`
	Explorer string `json:"explorer"`
	Version  int    `json:"version"`
}

//PeerResponseError struct to hold error response
type ConfigError struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}

//Error interface function
func (e ConfigError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v", e.Success, e.ErrorMessage)
}

func (s *ArkClient) AutoConfigureParams() {
	arkEnvParams := new(ArkEnvParams)
	configError := new(ConfigError)
	_, err := s.sling.New().Get("api/loader/autoconfigure").Receive(&arkEnvParams, configError)
	if err == nil {
		err = configError
	}

	_, err = s.sling.New().Get("api/blocks/getfees").Receive(&arkEnvParams, configError)
	if err == nil {
		err = configError
	}

}
