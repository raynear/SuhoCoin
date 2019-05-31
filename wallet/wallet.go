package wallet

import (
	"SuhoCoin/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

var Version string = "00"
var AddressChecksumLen int = 8

func NewWallet() *Wallet {
	pubkeyCurve := elliptic.P256() // P256이 가장 효율적이라 함 from https://safecurves.cr.yp.to

	privKey, e := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	err.ERR("Key generate Error", e)

	//	var publicKey ecdsa.PublicKey
	public := &privKey.PublicKey

	pubKey := append(public.X.Bytes(), public.Y.Bytes()...)

	wallet := Wallet{PrivateKey: *privKey, PublicKey: pubKey}
	return &wallet
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, e := RIPEMD160Hasher.Write(publicSHA256[:])
	err.ERR("RIPEMD160 Hash Error", e)
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:AddressChecksumLen]
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte(Version), pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := []byte(base58.Encode(fullPayload))

	return address
}
