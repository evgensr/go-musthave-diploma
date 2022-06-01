package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"log"
)

func GetHash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetShort(text string) string {
	// TODO: найти алгоритм который действительно будет предоставлять короткую ссылку
	return GetHash(text)
}

func GeneratorUUID() string {
	return uuid.New().String()
}

// Encrypted
func Encrypted(msg []byte, key string) ([]byte, error) {
	password := sha256.Sum256([]byte(key))
	aesblock, err := aes.NewCipher(password[:])
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	// создаём вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}

	dst := aesgcm.Seal(nil, nonce, msg, nil) // зашифровываем
	// добавим вектор
	buf, _ := hex.DecodeString(hex.EncodeToString(dst) + hex.EncodeToString(nonce))

	return buf, nil

}

func Decrypted(msg []byte, key string) ([]byte, error) {

	password := sha256.Sum256([]byte(key))

	aesblock, err := aes.NewCipher(password[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	// создаём вектор инициализации
	nonce := msg[len(msg)-aesgcm.NonceSize():]
	msgOriginal := msg[:len(msg)-aesgcm.NonceSize()]
	src2, err := aesgcm.Open(nil, nonce, msgOriginal, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	return src2, nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func AddSlash(s string) string {
	if len(s) < 1 {
		return ""
	}
	last := s[len(s)-1:]
	if last != `/` {
		s = s + `/`
	}
	return s
}
