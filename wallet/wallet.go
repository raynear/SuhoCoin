package wallet

import (
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

var Version []byte = []byte("00")
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

func NewWallet() string {
	pubkeyCurve := elliptic.P256() // P256이 가장 효율적이라 함 from https://safecurves.cr.yp.to

	privKey, e := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	util.ERR("Key generate Error", e)

	//	var publicKey ecdsa.PublicKey
	public := &privKey.PublicKey
	pubKey := append(public.X.Bytes(), public.Y.Bytes()...)

	wallet := Wallet{PrivateKey: *privKey, PublicKey: pubKey}

	address := wallet.GetAddress()
	fmt.Println("NewWallet")
	fmt.Println("Address : ", address, " is saved as file")

	fmt.Printf("Input Password: ")
	silentPassword, e := gopass.GetPasswdMasked()
	util.ERR("Password Input Error", e)

	EncryptedWallet := wallet.EncryptWallet(string(silentPassword))

	SaveToFile(address, EncryptedWallet)
	EncryptedWallet2 := LoadFromFile(address)
	NewWallet2 := DecryptWallet(EncryptedWallet2, string(silentPassword))

	fmt.Println("NewWallet2")
	fmt.Println("W PubKey2 : ", base58.Encode(NewWallet2.PublicKey))

	address2 := NewWallet2.GetAddress()

	fmt.Println("Address2 : ", address2, " is load from file")
	HashPubKey2 := HashPubKey(NewWallet2.PublicKey)
	fmt.Println("HashPubKey2 :", base58.Encode(HashPubKey2))

	return address
}

func (w *Wallet) EncodeWallet() []byte {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)

	e := encoder.Encode(w)
	util.ERR("wallet encode error", e)

	return content.Bytes()
}

func (w *Wallet) EncryptWallet(PassPhrase string) []byte {
	content := w.EncodeWallet()

	ciphertext := encrypt(content, PassPhrase)

	return ciphertext
}

func DecodeWallet(ByteWallet []byte) Wallet {
	var wallet Wallet
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewBuffer(ByteWallet))
	e := decoder.Decode(&wallet)
	util.ERR("decode wallet error", e)

	return wallet
}

func DecryptWallet(EncryptedWallet []byte, PassPhrase string) Wallet {
	plaintext := decrypt(EncryptedWallet, PassPhrase)

	NewWallet := DecodeWallet(plaintext)

	return NewWallet
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
	versionedPayload := append(Version, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)
	fmt.Println("W Address : ", address)

	return address
}

func GetPubKeyHashFromAddress(address string) []byte {
	fmt.Println("address : ", address)
	pubKeyHash := base58.Decode(address)
	pubKeyHash = pubKeyHash[2 : len(pubKeyHash)-4]
	fmt.Println("pubKeyHash : ", base58.Encode(pubKeyHash[:]))
	return pubKeyHash
}

func LoadFromFile(WalletFileName string) []byte {
	if _, e := os.Stat(WalletFileName + ".wallet"); os.IsNotExist(e) {
		util.ERR("WalletFile Not Exist", e)
	}

	fileContent, e := ioutil.ReadFile(WalletFileName + ".wallet")
	util.ERR("read wallet file error", e)

	return fileContent
}

func SaveToFile(WalletFileName string, EncryptedWallet []byte) {
	e := ioutil.WriteFile(WalletFileName+".wallet", EncryptedWallet, 0644)
	util.ERR("Write wallet file Error", e)
}
