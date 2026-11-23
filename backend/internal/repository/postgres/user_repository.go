package postgres

import (
	"context"
	"crypto/sha256"
	"ecommerce-backend/internal/domain"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Roles").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	normalized := strings.ToLower(strings.TrimSpace(email))
	err := r.db.WithContext(ctx).Preload("Roles").Where("LOWER(email) = ?", normalized).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) AssignRole(ctx context.Context, userID int64, roleID int) error {
	userRole := domain.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Where(domain.UserRole{UserID: userID, RoleID: roleID}).FirstOrCreate(&userRole).Error
}

func (r *UserRepository) CreateVerificationToken(ctx context.Context, token *domain.VerificationToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *UserRepository) GetVerificationToken(ctx context.Context, tokenStr string, tokenType string) (*domain.VerificationToken, error) {
	var token domain.VerificationToken
	hash := sha256.Sum256([]byte(tokenStr))
	tokenHash := hex.EncodeToString(hash[:])
	err := r.db.WithContext(ctx).
		Where("token = ? AND token_type = ? AND expires_at > NOW()", tokenHash, tokenType).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *UserRepository) DeleteVerificationToken(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.VerificationToken{}, id).Error
}

func (r *UserRepository) DeleteVerificationTokensByUser(ctx context.Context, userID int64, tokenType string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND token_type = ?", userID, tokenType).
		Delete(&domain.VerificationToken{}).Error
}

func (r *UserRepository) CreateAuthSession(ctx context.Context, session *domain.AuthSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *UserRepository) GetActiveAuthSession(ctx context.Context, tokenHash string) (*domain.AuthSession, error) {
	var session domain.AuthSession
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidSession
		}
		return nil, err
	}
	return &session, nil
}

func (r *UserRepository) RevokeAuthSession(ctx context.Context, tokenHash string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.AuthSession{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", now).Error
}

func (r *UserRepository) RevokeAllUserSessions(ctx context.Context, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.AuthSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}
