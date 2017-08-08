## How to install
```
$> go build
```

## How to run (windows)
```
$> ./arkgoserver.exe
```

## How to run (linux)
```
$> ./arkgoserver
```

## Run server in background
```
nohup ./arkgoserver &
```

## API's

* [http://localhost:54000/voters/rewards](http://localhost:54000/voters/rewards)
* [http://localhost:54000/voters/blocked](http://localhost:54000/voters/blocked)
* [http://localhost:54000/delegate](http://localhost:54000/delegate)
* [http://localhost:54000/delegate/config](http://localhost:54000/delegate/config)
* [http://localhost:54000/delegate/paymentruns](http://localhost:54000/delegate/paymentruns)
* [http://localhost:54000/delegate/paymentruns/details](http://localhost:54000/delegate/paymentruns/details)

## How to filter API

Specific Payment Run:
* http://localhost:54000/delegate/paymentruns/details?parentid=1

Specific Voter:
* http://localhost:54000/delegate/paymentruns/details?address=D5St8ot3asrxYW3o63EV3bM1VC6UBKMUfE

Specific Voter and Specific Run: 
* http://localhost:54000/delegate/paymentruns/details?parentid=1&address=D5St8ot3asrxYW3o63EV3bM1VC6UBKMUfE
