# ARKGO-POOL Documentation

This app was firstly developed to show the basic usage of ark-go library. 
The app functionalities are currently targeted at delegate functionalities - to setup automated script payments for delegate voters.

See the demo console app shell recording here:
https://asciinema.org/a/5yndxl794ncfpmjoqftuaiodm?t=8.

## How to setup a simple and efficient delegate payment system with arkgo-pool
It's  very straightforward *four(4) step process* as described bellow:
### 1. Download the binaries from the latest release
ARKGO-Pool is compiled into a small executable file that can be run on linux or windows. 

**No additional frameworks and installations needed. All you need to do is unzip the binary and run the application.**

**You can also run it from your secure USB stick.**
- [Download and unzip latest release](https://github.com/kristjank/ark-go/releases)

### 2. Edit you sharing settings in config.toml; to setup payment shares
Config.toml is a simple.ini setup file, where you specify the share percentages for your:
- voters,
- reserve fund,
- costs and
- personal fund. 

**It is very simple and intuitive.**
Just edit [settings/sample.config.toml](https://github.com/kristjank/ark-go/blob/master/cmd/arkgopool/settings/sample.config.toml) to your preferences and rename it to config.toml.

### 3. Run the executable to test the calculation
- $> ./arkgo-gui //Linux
- c:\ark-go-payment>arkgo-gui.exe  //Windows

See the **demo video** of console app shell recording here:
https://asciinema.org/a/5yndxl794ncfpmjoqftuaiodm?t=8.

### 4. When all is setup, run it in silent mode: 
To run the script in silent mode run:
**$>arkgo-gui -silent=true**

You can add it to you cronjob if you want to have automated payments.

**AND THAT'S ALL :)** 

### ARKGo-POOL functionalities include:
- Setting sharing properties in a simple config file: [settings/sample.config.toml](https://github.com/kristjank/ark-go/blob/master/cmd/arkgopool/settings/sample.config.toml). Copy file to config.toml for productive config. (cp sample.config.toml config.toml)
- Setting up delegate preferences in settings/config.toml. User can set his delegate address, sharing rations for voters and reward fund.
- Checking vote weight by voters, displaying fidelity in hoours and see the precalculated rewards before Sending.
- Manually sending the reward to users.
- Network switching (application can run on MAINNET as well as on DEVNET).
- Linking account for automated payments.
- Running in silent mode - to process automated reward payments to voters.
- Logging and overview of succesfull payments in a separate timestamped csv file, linked with sent and confirmed transactionsID from the blockchain (fields: Recepient, Amount, TimeStamp, TxID, VoteWieght).
- Database logging and REST server is included in the package
- Blocking of vote hopers

Application was developed to help testing of the core library and to provide a tool for delegates to process automated payments/ or just manual payments.

## Short guide for running from source

### Run from source
```
$> cd $GOPATH/src/github.com/kristjank/ark-go/cmd/arkgopool
$> go run arkgo-gui.go
```

### Build
```
$> cd $GOPATH/src/github.com/kristjank/ark-go/cmd/arkgopool
$> go build
```

### Run in Windows
```
$> cd $GOPATH/src/github.com/kristjank/ark-go/cmd/arkgopool
$> ./arkgo-gui.exe
```

### Run in Linux
```
$> cd $GOPATH/src/github.com/kristjank/ark-go/cmd/arkgopool
$> ./arkgopool
```

### Configuration

```
$> cd arkgopool/settings
$> cp sample.config.toml config.toml
$> nano config.toml
```

## Complete installation guide (ubuntu)

You will have to slightly adapt it to use the latest version. (Golang Download: https://golang.org/dl/)

```
$> curl -O https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
$> tar xvf go1.8.3.linux-amd64.tar.gz
$> sudo chown -R root:root ./go
$> sudo mv go /usr/local
$> sudo nano ~/.profile
```
Add at the end of the files:
```
export GOPATH=$HOME/work
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```
Save the files.
```
$> source ~/.profile
$> mkdir work
$> go get github.com/kristjank/ark-go
$> cd ~/work/src/github.com/kristjank/ark-go
$> go get ./...
$> cd ~/work/src/github.com/kristjank/ark-go/cmd/arkgopool
$> go build
$> cd
$> ln -s ~/work/src/github.com/kristjank/ark-go/cmd/arkgopool arkgopool
```
Now you can run with:
```
$> cd arkgopool
$> ./arkgopool
```

Create and edit config.toml to define your own settings.

```
$> cd arkgopool/settings
$> cp sample.config.toml config.toml
$> nano config.toml
```

## How to automate payment run

First, make sure that you did link your account. (Option 4)

Create a folder paymentlog in arkgopool:

```
$> cd arkgopool
```

Create such a file go-payment.sh in home:

```
#!/bin/bash
cd ~/work/src/github.com/kristjank/ark-go/main
./arkgopool -silent=true
```

Make it executable:

```
$> chmod +x go-payment.sh
```

Open crontab:

```
$> crontab -e
```

Add a line at the top (replace ~ with your absolute root path):

```
GOPATH=~/work
```

Add a line at the end and an empty line behind:

```
30 17 * * * ~/go-payment.sh
```

Save and now your payments will done automatic every day. Enjoy!


Thank you for checking in.

For more information about ark-go see the core package. That is where it all happens: [/core](/core)

## Authors
Chris (kristjan.kosic@gmail.com), with a lot of help from FX Thoorens fx@ark.io and ARK Community

## Contributors
Jarunik (jens@hviid.ch), with a lot of help from Chris

## Support this project
![Ark Logo](/raw/ark-logo-60x60.png)
Ark address:``AUgTuukcKeE4XFdzaK6rEHMD5FLmVBSmHk``


# License
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Copyright (c) 2017 ARK, chris

--chris
