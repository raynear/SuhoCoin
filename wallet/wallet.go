package wallet

import (
	"SuhoCoin/config"
	"SuhoCoin/util"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"github.com/howeyc/gopass"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

var Version string = "0"
var AddressChecksumLen int = 4

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, e := cipher.NewGCM(block)
	util.ERR("aes cipher gcm error", e)
	nonce := make([]byte, gcm.NonceSize())
	if _, e = io.ReadFull(rand.Reader, nonce); e != nil {
		panic(e.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, e := aes.NewCipher(key)
	util.ERR("aes cipher error", e)
	gcm, e := cipher.NewGCM(block)
	util.ERR("aes cipher gcm error", e)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, e := gcm.Open(nil, nonce, ciphertext, nil)
	util.ERR("Decrypt Error(Maybe Wrong Password)", e)
	if e != nil {
		return []byte{}
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}

func NewWallet() *Wallet {
	pubkeyCurve := elliptic.P256() // P256이 가장 효율적이라 함 from https://safecurves.cr.yp.to

	privKey, e := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	util.ERR("Key generate Error", e)

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
	util.ERR("RIPEMD160 Hash Error", e)
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:AddressChecksumLen]
}

func (w Wallet) GetAddress() string {
	fmt.Println("W PubKey : ", base58.Encode(w.PublicKey))
	pubKeyHash := HashPubKey(w.PublicKey)
	fmt.Println("W PubKeyHash : ", base58.Encode(pubKeyHash))
	versionedPayload := append([]byte(Version), pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)
	fmt.Println("W Address : ", address)

	return address
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	e := wallets.LoadFromFile()

	return &wallets, e
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Println("newAddress : ", address)

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile() error {
	if _, e := os.Stat(config.V.GetString("WalletFile")); os.IsNotExist(e) {
		return e
	}

	fileContent, e := ioutil.ReadFile(config.V.GetString("WalletFile"))
	util.ERR("read wallet file error", e)

	fmt.Printf("Input Password: ")
	silentPassword, e := gopass.GetPasswdMasked()
	util.ERR("Password Input Error", e)

	plaintext := decrypt(fileContent, string(silentPassword))

	if bytes.Compare(plaintext, []byte("")) == 0 {
		fmt.Println("Wrong Password")
		return nil
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(plaintext))
	e = decoder.Decode(&wallets)
	util.ERR("decode wallet error", e)

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws Wallets) SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)

	e := encoder.Encode(ws)
	util.ERR("wallet encode error", e)

	fmt.Printf("Input Password: ")
	silentPassword, e := gopass.GetPasswdMasked()
	util.ERR("Password Input Error", e)

	ciphertext := encrypt(content.Bytes(), string(silentPassword))

	e = ioutil.WriteFile(config.V.GetString("WalletFile"), ciphertext, 0644)
	util.ERR("Write wallet file Error", e)
}
