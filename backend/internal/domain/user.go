package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string         `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Email      string         `json:"email" gorm:"column:email;type:varchar(100);unique;not null"`
	Password   string         `json:"-" gorm:"column:password;type:varchar(255);not null"`
	Phone      string         `json:"phone" gorm:"column:phone;type:varchar(20)"`
	Address    string         `json:"address" gorm:"column:address;type:text"`
	Photo      string         `json:"photo" gorm:"column:photo;type:varchar(255)"`
	Lat        string         `json:"lat" gorm:"column:lat;type:varchar(50)"`
	Lng        string         `json:"lng" gorm:"column:lng;type:varchar(50)"`
	IsVerified bool           `json:"is_verified" gorm:"column:is_verified;default:false"`
	Roles      []Role         `json:"roles,omitempty" gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

type Role struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name;type:varchar(50);unique;not null"`
}

func (Role) TableName() string {
	return "roles"
}

type UserRole struct {
	ID     int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int64 `json:"user_id" gorm:"column:user_id;not null"`
	RoleID int   `json:"role_id" gorm:"column:role_id;not null"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type VerificationToken struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64          `json:"user_id" gorm:"column:user_id;not null"`
	Token     string         `json:"token" gorm:"column:token;type:varchar(255);not null"`
	TokenType string         `json:"token_type" gorm:"column:token_type;type:varchar(50);not null"` // EMAIL_VERIFICATION, RESET_PASSWORD
	ExpiresAt time.Time      `json:"expires_at" gorm:"column:expires_at;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (VerificationToken) TableName() string {
	return "verification_tokens"
}

// DTOs
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=50"`
	Phone    string `json:"phone" validate:"omitempty,max=20"`
	Address  string `json:"address" validate:"omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	AssignRole(ctx context.Context, userID int64, roleID int) error
	CreateVerificationToken(ctx context.Context, token *VerificationToken) error
	GetVerificationToken(ctx context.Context, token string, tokenType string) (*VerificationToken, error)
	DeleteVerificationToken(ctx context.Context, id int64) error
}

type AuthUseCase interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	GetProfile(ctx context.Context, userID int64) (*User, error)
	VerifyEmail(ctx context.Context, token string) error
}
