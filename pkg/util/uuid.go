package util

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID 生成一个新的UUID并返回其字符串表示
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUUIDWithoutHyphen 生成一个没有连字符的UUID字符串
func GenerateUUIDWithoutHyphen() string {
	return strings.ReplaceAll(GenerateUUID(), "-", "")
}

// ParseUUID 将字符串解析为UUID对象
// 如果解析失败，返回错误
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// IsValidUUID 检查字符串是否是有效的UUID格式
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
