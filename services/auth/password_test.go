package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password %v", err)
	}
	
	if hash == "" {
		t.Errorf("expected hash shouldn't empty")
	}
	
	if hash == "password" {
		t.Errorf("expected hash should be different from password")
	}
}

func TestCompareHashAndPassword(t *testing.T)  {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password %v", err)
	}
	if !ComparePasswords(hash, []byte("password")){
		t.Errorf("expected hashed password to match the hash")
	}
	if ComparePasswords(hash, []byte("notpassword")){
		t.Errorf("expected hashed password to not match the hash")
	}
}