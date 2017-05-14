package arkcoin

var (
	//BitcoinMain is params for main net.
	ArkCoinMain = &Params{
		DumpedPrivateKeyHeader: []byte{170}, //wif
		AddressHeader:          23,
		P2SHHeader:             5,
		HDPrivateKeyID:         []byte{0x04, 0x88, 0xad, 0xe4},
		HDPublicKeyID:          []byte{0x04, 0x88, 0xb2, 0x1e},
	}
	BitcoinMain = &Params{
		DumpedPrivateKeyHeader: []byte{128},
		AddressHeader:          0,
		P2SHHeader:             5,
		HDPrivateKeyID:         []byte{0x04, 0x88, 0xad, 0xe4},
		HDPublicKeyID:          []byte{0x04, 0x88, 0xb2, 0x1e},
	}
	//BitcoinTest is params for test net.
	BitcoinTest = &Params{
		DumpedPrivateKeyHeader: []byte{239},
		AddressHeader:          111,
		P2SHHeader:             196,
		HDPrivateKeyID:         []byte{0x04, 0x35, 0x83, 0x94},
		HDPublicKeyID:          []byte{0x04, 0x35, 0x87, 0xcf},
	}
	//MonacoinMain is params for monacoin main net.
	MonacoinMain = &Params{
		DumpedPrivateKeyHeader: []byte{178, 176},
		AddressHeader:          50,
		P2SHHeader:             5,
		HDPrivateKeyID:         []byte{0x04, 0x88, 0xad, 0xe4},
		HDPublicKeyID:          []byte{0x04, 0x88, 0xb2, 0x1e},
	}
)
