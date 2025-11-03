package util

import "testing"

func TestEncryptAESGCM(t *testing.T) {
	plaintext := []byte("Hello, World!aaaaadfadfadfadfadfadfa")
	key := []byte(AESGCMDefaultKey)
	encrypted, err := EncryptAESGCM(plaintext, key)
	if err != nil {
		t.Errorf("EncryptAESGCM failed: %v", err)
		return
	}
	t.Logf("Encrypted data: %s", encrypted)
	decrypted, err := DecryptAESGCM(encrypted, key)
	if err != nil {
		t.Errorf("DecryptAESGCM failed: %v", err)
		return
	}
	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted data does not match original data")
	}
}
