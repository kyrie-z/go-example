package u_rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// generate
func GenRsaKey(bits int) (prikey, pubkey []byte) {
	keyPair, _ := rsa.GenerateKey(rand.Reader, bits)
	derPrikey := x509.MarshalPKCS1PrivateKey(keyPair)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrikey,
	}
	prikey = pem.EncodeToMemory(block)
	derPubkey, _ := x509.MarshalPKIXPublicKey(&keyPair.PublicKey)

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPubkey,
	}
	pubkey = pem.EncodeToMemory(block)
	return prikey, pubkey
}

// save key
func SaveKey(key []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	block, _ := pem.Decode(key)
	err = pem.Encode(f, block)
	if err != nil {
		return err
	}
	return nil
}

// sign
func RsaSignWithSha256(plain []byte, prikey *rsa.PrivateKey) []byte {
	//hash
	h := sha256.New()
	h.Write(plain)
	hashData := h.Sum(nil)
	// fmt.Println(hashData)
	//rsa sign
	signData, err := rsa.SignPKCS1v15(rand.Reader, prikey, crypto.SHA256, hashData)
	if err != nil {
		fmt.Println("sign fail!!")
		return nil
	}
	return signData
}

// verify
func RsaVerifyWithSha256(plain, signData []byte, pubkey *rsa.PublicKey) bool {
	// plain hash
	hashData := sha256.Sum256(plain)
	//verify
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashData[:], signData)
	if err != nil {
		fmt.Println("verify fail!!")
		return false
	}
	return true
}

// encrypt
func RsaEncrypt(plain []byte, pubkey *rsa.PublicKey) []byte {
	cipherData, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, plain)
	if err != nil {
		fmt.Println("encrypt fail!!")
		return nil
	}
	return cipherData
}

// decrypt
func RsaDecrypt(cipherData []byte, pribkey *rsa.PrivateKey) []byte {
	plain, err := rsa.DecryptPKCS1v15(rand.Reader, pribkey, cipherData)
	if err != nil {
		fmt.Println("decrypt fail!!")
		return nil
	}
	return plain

}
