package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	// 定义可用字符集（0-9A-Za-z）
	randomStringChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// GenerateRandomString 生成指定长度的随机字符串（包含0-9、A-Z、a-z）
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive integer")
	}

	// 计算需要的随机字节数（每个字符需要log62(256) ≈ 1.3字节）
	// 取2倍长度保证足够随机性
	byteLen := length * 2
	b := make([]byte, byteLen)
	
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// 将随机字节转换为十六进制字符串
	hexStr := hex.EncodeToString(b)
	
	// 转换字符集
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		// 取两个十六进制字符（0-255）
		num := hexStr[i*2 : (i+1)*2]
		val, err := hex.DecodeString(num)
		if err != nil {
			return "", fmt.Errorf("hex conversion failed: %w", err)
		}
		
		// 映射到62个字符
		result[i] = randomStringChars[int(val[0])%len(randomStringChars)]
	}

	return string(result), nil
}