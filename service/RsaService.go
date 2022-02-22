package service

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"licenseGenerator/dao"
	"licenseGenerator/models"
	"log"
)

func GenerateKey() models.RsaKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	publicKey := &privateKey.PublicKey

	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	X509PublicKey := x509.MarshalPKCS1PublicKey(publicKey)

	privateKeyString := base64.StdEncoding.EncodeToString(X509PrivateKey)
	publicKeyString := base64.StdEncoding.EncodeToString(X509PublicKey)

	var key = models.RsaKey{
		PublicKey:  publicKeyString,
		PrivateKey: privateKeyString,
	}
	err := dao.Db.Create(&key).Error
	if err != nil {
		panic(err)
	}
	return key
}

func DecodePrivateKeyString(privateKeyString string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func DecodePublicKeyString(publicKeyString string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyString)
	if err != nil {
		return nil, err
	}
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func Sign(privateKey *rsa.PrivateKey, origData []byte) (string, error) {
	msgHash := sha256.New()
	msgHash.Write(origData)
	msgHashSum := msgHash.Sum(nil)
	signatureBytes, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signatureBytes), nil
}

func Verify(publicKey *rsa.PublicKey, origData []byte, signatureString string) error {
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureString)
	if err != nil {
		return err
	}
	msgHash := sha256.New()
	msgHash.Write(origData)
	msgHashSum := msgHash.Sum(nil)

	err = rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSum, signatureBytes, nil)
	if err != nil {
		return err
	}
	return nil
}

// 公钥加密
func Encrypt(publicKey *rsa.PublicKey, origData []byte) (string, error) {
	partLen := publicKey.N.BitLen()/16 - 66
	chunks := Split(origData, partLen)
	var encrypts bytes.Buffer
	for _, chunk := range chunks {
		encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, chunk, []byte{})
		if err != nil {
			log.Println(len(chunk), publicKey.Size())
			return "", err
		}
		encrypts.Write(encrypted)
	}
	return base64.StdEncoding.EncodeToString(encrypts.Bytes()), nil
}

// 私钥解密
func Decrypt(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey, encryptString string) ([]byte, error) {
	partLen := publicKey.N.BitLen() / 16
	encryptBytes, err := base64.StdEncoding.DecodeString(encryptString)
	if err != nil {
		return nil, err
	}
	chunks := Split(encryptBytes, partLen)
	var decryptBytes bytes.Buffer
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, chunk, []byte{})
		if err != nil {
			return nil, err
		}
		decryptBytes.Write(decrypted)
	}
	return decryptBytes.Bytes(), err
}

// 数据分块
func Split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
