## ARKGO Load Testing script
This is a work in progress....

ArkGo Tester is a configurable executable script to test ARK BC performance and stability. 
Script is adjustable for payload size, iterations, multibroadcast, single node testing, etc... So you can simulate various situations...

## Set your configuration
Edit cfg/config.toml and start running tests.
You can set: 

```
[env]
txPerPayload=250                   #nr of transactions per payload
txIterations=1                     #nr of iterations of payload deliveries... all txIterations * txPerPayload = all tx being sent
txMultiBroadCast = 1               #nr of peers to multibroadcast to
txDescription="ARK-GO Testing program running" #txDescription
dbFileName = "db/testlog.db"

[account]
passphrase=""                                  #passphrase of the test tx sending account
secondPassphrase =""
recepient=""                                   #recepient of the transactions
```

## How to run...
Just run the executable file attached in the package...

## TODO
- different test type implementations 
  - same payload, one peer
  - same payload, multiple peers
  - mixed payload, multiple peers
  
  
--delegate chris
