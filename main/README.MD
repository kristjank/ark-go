# ARK-GO Client Demo Usage

This app is just to show the basic usage of ark-go library. The app functionalities are currently targeted at delegate functionalities, but from code various API calls can be seen. Further functionalities will be added based on community feature requests.

To see all the calls and usage of ARK-GO client library check the name_test.go classes in the core package!

ArkGo-Gui functionalities include:
- Setting up delegate preferences in [settings/config.toml](https://github.com/kristjank/ark-go/blob/master/main/settings/config.toml). User can set his delegate address, sharing rations for voters and reward fund.
- Checking vote weight by voters, displaying fidelity in hoours and see the precalculated rewards before Sending.
- Manually sending the reward to users.
- Network switching (application can run on MAINNET as well as on DEVNET).
- Linking account for automated payments.
- Running in silent mode - to process automated reward payments to voters.
- Logging and overview of succesfull payments in a separate timestamped csv file, linked with sent and confirmed transactionsID from the blockchain (fields: Recepient, Amount, TimeStamp, TxID, VoteWieght).

Application was developed to help testing of the core library and to provide a tool for delegates to process automated payments/ or just manual payments.

Thank you for checking in. 

For more information about ark-go see the core package. That is where it all happens: https://github.com/kristjank/ark-go/tree/master/core

--chris
