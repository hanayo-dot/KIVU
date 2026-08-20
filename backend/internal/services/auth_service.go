package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"github.com/hanayo-dot/KIVU/backend/config"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// AuthService manages user registration, login, JWT token generation, and refresh token rotation.
type AuthService struct {
	cfg *config.Config
	db  *sqlx.DB
}

// NewAuthService initializes AuthService.
func NewAuthService(cfg *config.Config, db *sqlx.DB) *AuthService {
	return &AuthService{
		cfg: cfg,
		db:  db,
	}
}

// RegisterFarmer registers a new farmer and returns access & refresh tokens.
func (s *AuthService) RegisterFarmer(req models.RegisterRequest) (*models.TokenResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	location := req.LocationName
	if location == "" {
		location = "Lake Victoria Basin"
	}

	var farmer models.Farmer
	query := `
		INSERT INTO farmers (name, phone_number, password_hash, location_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, phone_number, location_name, created_at`

	err = s.db.Get(&farmer, query, req.Name, req.PhoneNumber, string(hashedPassword), location)
	if err != nil {
		return nil, fmt.Errorf("farmer registration failed (phone number may already exist): %w", err)
	}

	return s.issueTokens(farmer)
}

// LoginFarmer verifies credentials and returns a new access & refresh token pair.
func (s *AuthService) LoginFarmer(req models.LoginRequest) (*models.TokenResponse, error) {
	var farmer models.Farmer
	err := s.db.Get(&farmer, `SELECT id, name, phone_number, password_hash, location_name, created_at FROM farmers WHERE phone_number = $1`, req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(farmer.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid phone number or password")
	}

	return s.issueTokens(farmer)
}

// RefreshToken rotates refresh tokens and issues a new access token.
func (s *AuthService) RefreshToken(rawRefreshToken string) (*models.TokenResponse, error) {
	tokenHash := s.hashToken(rawRefreshToken)

	var tokenRecord models.RefreshToken
	err := s.db.Get(&tokenRecord,
		`SELECT id, farmer_id, expires_at, created_at, revoked 
		 FROM refresh_tokens 
		 WHERE token_hash = $1`, tokenHash)
	if err != nil || tokenRecord.Revoked || time.Now().After(tokenRecord.ExpiresAt) {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Revoke old refresh token (Token Rotation)
	_, _ = s.db.Exec(`UPDATE refresh_tokens SET revoked = TRUE WHERE id = $1`, tokenRecord.ID)

	var farmer models.Farmer
	err = s.db.Get(&farmer, `SELECT id, name, phone_number, location_name, created_at FROM farmers WHERE id = $1`, tokenRecord.FarmerID)
	if err != nil {
		return nil, fmt.Errorf("farmer not found for token: %w", err)
	}

	return s.issueTokens(farmer)
}

// LogoutRevokeToken revokes the specified refresh token.
func (s *AuthService) LogoutRevokeToken(rawRefreshToken string) error {
	tokenHash := s.hashToken(rawRefreshToken)
	_, err := s.db.Exec(`UPDATE refresh_tokens SET revoked = TRUE WHERE token_hash = $1`, tokenHash)
	return err
}

func (s *AuthService) issueTokens(farmer models.Farmer) (*models.TokenResponse, error) {
	// 1. Generate JWT Access Token (Short-lived)
	accessExpDuration := time.Duration(s.cfg.JWTAccessExpiryMinutes) * time.Minute
	accessExpiresAt := time.Now().Add(accessExpDuration)

	claims := jwt.MapClaims{
		"farmer_id": farmer.ID.String(),
		"iat":       time.Now().Unix(),
		"exp":       accessExpiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHMAC256, claims)
	accessTokenStr, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// 2. Generate Opaque Refresh Token (Long-lived)
	rawRefreshToken, err := s.generateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshExpDuration := time.Duration(s.cfg.JWTRefreshExpiryDays) * 24 * time.Hour
	refreshExpiresAt := time.Now().Add(refreshExpDuration)
	tokenHash := s.hashToken(rawRefreshToken)

	// Save hashed refresh token to database
	_, err = s.db.Exec(
		`INSERT INTO refresh_tokens (farmer_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		farmer.ID, tokenHash, refreshExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpDuration.Seconds()),
		Farmer:       farmer,
	}, nil
}

func (s *AuthService) generateRandomToken(bytesLen int) (string, error) {
	b := make([]byte, bytesLen)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *AuthService) hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
