package util

import (
	"strings"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	if len(uuid) != 36 {
		t.Errorf("生成的UUID长度应为36，实际为%d", len(uuid))
	}
	if strings.Count(uuid, "-") != 4 {
		t.Errorf("UUID应该包含4个连字符，实际包含%d个", strings.Count(uuid, "-"))
	}
}

func TestGenerateUUIDWithoutHyphen(t *testing.T) {
	uuid := GenerateUUIDWithoutHyphen()
	if len(uuid) != 32 {
		t.Errorf("生成的无连字符UUID长度应为32，实际为%d", len(uuid))
	}
	if strings.Contains(uuid, "-") {
		t.Error("无连字符UUID不应包含连字符")
	}
}

func TestParseUUID(t *testing.T) {
	// 测试有效的UUID
	validUUID := GenerateUUID()
	_, err := ParseUUID(validUUID)
	if err != nil {
		t.Errorf("解析有效UUID失败: %v", err)
	}

	// 测试无效的UUID
	invalidUUID := "invalid-uuid"
	_, err = ParseUUID(invalidUUID)
	if err == nil {
		t.Error("解析无效UUID应该返回错误")
	}
}

func TestIsValidUUID(t *testing.T) {
	// 测试有效的UUID
	validUUID := GenerateUUID()
	if !IsValidUUID(validUUID) {
		t.Error("有效的UUID被判定为无效")
	}

	// 测试无效的UUID
	invalidCases := []string{
		"invalid-uuid",
		"123",
		"",
		"550e8400-e29b-41d4-a716-446655440000x",
	}

	for _, invalidUUID := range invalidCases {
		if IsValidUUID(invalidUUID) {
			t.Errorf("无效的UUID '%s' 被判定为有效", invalidUUID)
		}
	}
}
