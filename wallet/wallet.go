package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const (
	checksumlength = 4
	version = byte(0x00)
)
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func (w *Wallet) Address()[]byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionHash := append([]byte{version}, pubHash...)
	checkSum := CheckSum(versionHash)
	hash := append(versionHash, checkSum...)
	address := Base58Encode(hash)

	//fmt.Printf("pub key: %x\n", w.PublicKey)
	//fmt.Printf("pub hash: %x\n", pubHash)
	//fmt.Printf("address: %x\n", address)
	return address
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	realCheckSum := pubKeyHash[len(pubKeyHash) - checksumlength:]
	version := pubKeyHash[0]
	targetCheckSum := CheckSum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(realCheckSum, targetCheckSum) == 0
}
func NewKeyPair() (ecdsa.PrivateKey, []byte){
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func PublicKeyHash(pubkey []byte) []byte {
	pubHash := sha256.Sum256(pubkey)
	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipMD := hasher.Sum(nil)
	return publicRipMD
}

func CheckSum (payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:checksumlength]
}
