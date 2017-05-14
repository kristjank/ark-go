package arkcoin

//BitcoinMain is params for main net.
var (
	ArkCoinMain = &Params{
		DumpedPrivateKeyHeader: []byte{170}, //wif
		AddressHeader:          23,          //0x17
		//P2SHHeader:             5,
		//HDPrivateKeyID:         []byte{0x04, 0x88, 0xad, 0xe4},
		//HDPublicKeyID:          []byte{0x04, 0x88, 0xb2, 0x1e},
	}
	BitcoinMain = &Params{
		DumpedPrivateKeyHeader: []byte{128},
		AddressHeader:          0,
		//P2SHHeader:             5,
		//HDPrivateKeyID:         []byte{0x04, 0x88, 0xad, 0xe4},
		//HDPublicKeyID:          []byte{0x04, 0x88, 0xb2, 0x1e},
	}
	//BitcoinTest is params for test net.
	BitcoinTest = &Params{
		DumpedPrivateKeyHeader: []byte{239},
		AddressHeader:          111,
		P2SHHeader:             196,
		//HDPrivateKeyID:         []byte{0x04, 0x35, 0x83, 0x94},
		//HDPublicKeyID:          []byte{0x04, 0x35, 0x87, 0xcf},
	}
)
