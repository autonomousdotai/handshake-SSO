package utils

import (
    "io"
    "fmt"
    "errors"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "github.com/segmentio/ksuid"
)

func HashNewUID(key string) (string, string, error) {
    bkey := []byte(key)
    uuid := ksuid.New().String()
    passpharse, err := HashEncrypt(bkey, uuid)
    if err != nil {
        return "", "", err
    }
    return uuid, passpharse, nil
}

func HashEncrypt(key []byte, text string) (string, error) {
    plaintext := []byte(text)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize + len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func HashDecrypt(key []byte, cryptoText string) (string, error) {
    ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    if len(ciphertext) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)

    stream.XORKeyStream(ciphertext, ciphertext)

    return fmt.Sprintf("%s", ciphertext), nil
}

