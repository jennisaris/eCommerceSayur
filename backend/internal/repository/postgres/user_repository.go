package postgres

import (
	"context"
	"ecommerce-backend/internal/domain"
	"errors"

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
	err := r.db.WithContext(ctx).Preload("Roles").Where("email = ?", email).First(&user).Error
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
	err := r.db.WithContext(ctx).
		Where("token = ? AND token_type = ? AND expires_at > NOW()", tokenStr, tokenType).
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
