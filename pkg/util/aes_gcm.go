package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

const (
	AESGCMDefaultKey = "Adf2eb05cf39a0994c4f200d8b5af901"
)

// EncryptAESGCM 使用 AES-256-GCM 加密数据
// 参数：
//
//	plaintext - 要加密的原文
//	key       - 32字节密钥（AES-256）
//
// 返回：base64(nonce + ciphertext)
func EncryptAESGCM(plaintext, key []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce（每次加密都不同）
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密（并自动添加认证标签）
	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	// 拼接 nonce + ciphertext 一起输出
	result := append(nonce, ciphertext...)

	// 转为 base64 字符串便于存储或传输
	return base64.StdEncoding.EncodeToString(result), nil
}

// DecryptAESGCM 解密 base64 编码的 AES-GCM 密文
func DecryptAESGCM(encoded string, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	// 解密并校验
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GenerateRandomKey 生成一个随机 AES-256 密钥（32字节）
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}
