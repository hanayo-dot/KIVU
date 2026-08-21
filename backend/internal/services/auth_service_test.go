package services

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHashingAndVerification(t *testing.T) {
	password := "SecretDemoPass123!"

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Compare correct password
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Errorf("expected password comparison to succeed, got error: %v", err)
	}

	// Compare wrong password
	err = bcrypt.CompareHashAndPassword(hashed, []byte("WrongPass"))
	if err == nil {
		t.Errorf("expected password comparison to fail for wrong password")
	}
}

func TestJWTTokenClaimsGenerationAndParsing(t *testing.T) {
	secretKey := []byte("test_jwt_secret_key_32bytes_len!")
	farmerID := uuid.New()
	phone := "+254712345678"

	// Create JWT Claims
	claims := jwt.MapClaims{
		"farmer_id":    farmerID.String(),
		"phone_number": phone,
		"exp":          time.Now().Add(1 * time.Hour).Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	// Parse Token back
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		t.Fatalf("expected token to be valid, got err: %v", err)
	}

	parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected MapClaims type")
	}

	if parsedClaims["farmer_id"] != farmerID.String() {
		t.Errorf("expected farmer_id %s, got %v", farmerID.String(), parsedClaims["farmer_id"])
	}
	if parsedClaims["phone_number"] != phone {
		t.Errorf("expected phone_number %s, got %v", phone, parsedClaims["phone_number"])
	}
}

func TestRefreshTokenSHA256Hashing(t *testing.T) {
	rawToken := "5a9d8f7e6c5b4a3f2e1d"
	hashBytes := sha256.Sum256([]byte(rawToken))
	expectedHash := hex.EncodeToString(hashBytes[:])

	if len(expectedHash) != 64 {
		t.Errorf("expected 64 character hex string for SHA-256, got length %d", len(expectedHash))
	}
}
