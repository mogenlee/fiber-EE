package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "test123456"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}

	if hash == password {
		t.Fatal("Hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "test123456"

	hash, _ := HashPassword(password)

	// 正确密码
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}

	// 错误密码
	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
}
