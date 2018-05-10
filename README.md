![ARK-Go](https://github.com/kristjank/ark-go/blob/master/raw/Ark-GO.png)

[![GitHub issues](https://img.shields.io/github/issues/kristjank/ark-net.svg)](https://github.com/kristjank/ark-go/issues)&nbsp;[![GitHub forks](https://img.shields.io/github/forks/kristjank/ark-net.svg)](https://github.com/kristjank/ark-go/network)&nbsp;[![GitHub stars](https://img.shields.io/github/stars/kristjank/ark-net.svg)](https://github.com/kristjank/ark-go/stargazers)&nbsp;[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kristjank/ark-go/master/LICENSE)

## Overview

A short intro about the different directories within this repository:

* [/arkcoin](/arkcoin) - Coin configuration and helper functions
* [/core](/core) - API to use the Ark blockchain out of Golang
* [/cmd/arkgopool](/cmd/arkgopool) - Client for delegate profit sharing pools
* [/cmd/arkgoserver](/cmd/arkgoserver) - Server for delgates profit sharing pools
* [/raw](/raw) - Images and other raw files

You can find more information about the different functionalities in their folder.

## Why Ark-GO
GoLang is an open source programming language developed by Google and designed for building fast, simple and reliable software. It is not about theoretical concepts such as monads and virtual inheritance, but more about **hands-on experience**.

Ark-GO is the ARK Ecosystem library client implemented in GOLANG programming language. It implements all most relevant ARK functionalities to help you  **develop efficient, fast and scalable GOLANG applications built upon ARK platform**. It provides also low level access to ARK so you can easily build your application on top of it.

## A library demo app: ARKGO-GUI
A demo client app was developed to test the goark library package dependencies. More about the demo gui client: [/cmd/arkgopool](/cmd/arkgopool). It targets delegate and basic account functionalites, plus silent mode - to be able to run automated reward payments.

See the demo console app shell recording here:
https://asciinema.org/a/5yndxl794ncfpmjoqftuaiodm?t=8.

## How to install?
```
$> go get github.com/kristjank/ark-go
```

## How to get dependencies?
```
$> go get ./...
```

## How to get started?
All ark-node services have available reponses have their struct representations. It's best to let the code do the speaking. Every class implementation has it's own test class. **So it's best to start learning by looking at actual test code**.

## Ark-GO Client Usage
**First call should be network selection, so all settings can initialize from the peers before going into action.**  By default MAINNET is active.

### Init
[GoDoc documentation available on this link](https://godoc.org/github.com/kristjank/ark-go/core)
```go
import "github.com/kristjank/ark-go/core"
var arkclient = core.NewArkClient(nil)
```

### Communication
Queries to the blockchain are done with the Query struct parameters:

```go
params := TransactionQueryParams{Limit: 10, SenderID: senderID}
```
... and the results -  reponse is also parametrized.
```go
transResponse, _, err := arkapi.ListTransaction(params)
if transResponse.Success {
		log.Println(t.Name(), "Success, returned", transResponse.Count, "transactions")
	} else {
		t.Error(err.Error())
	}
```

### Other call samples
```go
//usage samples
deleResp, _, _ := arkclient.GetDelegate(params)

//switch networks
arkclient = arkclient.SetActiveConfiguration(core.DEVNET) //or core.MAINNET
//create and send tx
arkapi := NewArkClient(nil)
recepient := "address"
passphrase := "pass"

tx := CreateTransaction(recepient,1,"ARK-GOLang is saying whoop whooop",passphrase, "")
payload.Transactions = append(payload.Transactions, tx)
res, httpresponse, err := arkapi.PostTransaction(payload)
```
## More information about ARK Ecosystem and etc
* [ARK Ecosystem Wiki](https://github.com/ArkEcosystem/wiki)

Please, use github issues for questions or feedback. For confidential requests or specific demands, contact us on our public channels.

## Authors
Chris (kristjan.kosic@gmail.com), with a lot of help from FX Thoorens fx@ark.io and ARK Community

Jarunik (jens@hviid.ch), with a lot of help from Chris

## Support this project
![Ark Logo](/raw/ark-logo-60x60.png)
Ark address:``AUgTuukcKeE4XFdzaK6rEHMD5FLmVBSmHk``


# License
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Copyright (c) 2017 ARK, chris
