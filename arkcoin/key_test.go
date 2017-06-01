package arkcoin

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/btcsuite/btcd/btcec"
)

func TestKeys2(t *testing.T) {
	key, err := Generate(ArkCoinMain)
	if err != nil {
		t.Errorf(err.Error())
	}
	adr := key.PublicKey.Address()
	log.Println("address=", adr)
	wif := key.WIFAddress()
	log.Println("wif=", wif)

	key2, err := FromWIF(wif, ArkCoinMain)
	if err != nil {
		t.Errorf(err.Error())
	}
	adr2 := key2.PublicKey.Address()
	log.Println("address2=", adr2)

	if adr != adr2 {
		t.Errorf("key unmatched")
	}
}

func TestKeys(t *testing.T) {
	key, err := Generate(ArkCoinMain)
	if err != nil {
		t.Errorf(err.Error())
	}
	adr := key.PublicKey.Address()
	log.Println("address=", adr)
	wif := key.WIFAddress()
	log.Println("wif=", wif)

	key2, err := FromWIF(wif, ArkCoinMain)
	if err != nil {
		t.Errorf(err.Error())
	}
	adr2 := key2.PublicKey.Address()
	log.Println("address2=", adr2)

	if adr != adr2 {
		t.Errorf("key unmatched")
	}

}

func TestSign(t *testing.T) {
	seed := make([]byte, 32)
	_, err := hex.Decode(seed, []byte("3954e0c9a3ce58a8dca793e214232e569ff0cb9da79689ca56d0af614227d540"))
	if err != nil {
		t.Fatal(err)
	}
	s256 := btcec.S256()
	priv, pub := btcec.PrivKeyFromBytes(s256, seed)
	public := PublicKey{
		PublicKey:    pub,
		isCompressed: false,
	}
	private := PrivateKey{
		PrivateKey: priv,
		PublicKey:  &public,
	}
	data := []byte("test data")
	sig, err := private.Sign(data)
	if err != nil {
		t.Fatal(err)
	}
	if err = private.PublicKey.Verify(sig, data); err != nil {
		t.Error(err)
	}
	data2 := []byte("invalid test data")
	if err = private.PublicKey.Verify(sig, data2); err == nil {
		t.Error("cannot verify")
	}
	log.Println(err)
}
