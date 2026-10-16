package usecase

import (
	"context"
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/hash"
	"ecommerce-backend/internal/pkg/jwt"
	"errors"
	"time"
)

type AuthUseCase struct {
	userRepo domain.UserRepository
	cartRepo domain.CartRepository
	jwtSvc   *jwt.JWTService
}

func NewAuthUseCase(userRepo domain.UserRepository, cartRepo domain.CartRepository, jwtSvc *jwt.JWTService) domain.AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		cartRepo: cartRepo,
		jwtSvc:   jwtSvc,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if email already taken
	existing, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Phone:      req.Phone,
		Address:    req.Address,
		IsVerified: true, // auto verify for seamless onboarding
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Assign 'Customer' role (role_id = 2)
	_ = u.userRepo.AssignRole(ctx, newUser.ID, 2)

	// Pre-create empty cart for user
	_ = u.cartRepo.Create(ctx, &domain.Cart{UserID: newUser.ID})

	// Fetch fresh user with roles
	userWithRoles, err := u.userRepo.GetByID(ctx, newUser.ID)
	if err != nil {
		userWithRoles = newUser
	}

	var roleNames []string
	for _, r := range userWithRoles.Roles {
		roleNames = append(roleNames, r.Name)
	}
	if len(roleNames) == 0 {
		roleNames = []string{"Customer"}
	}

	token, err := u.jwtSvc.GenerateToken(newUser.ID, newUser.Email, roleNames)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  userWithRoles,
	}, nil
}

func (u *AuthUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	var roleNames []string
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Name)
	}

	token, err := u.jwtSvc.GenerateToken(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (u *AuthUseCase) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, userID)
}

func (u *AuthUseCase) VerifyEmail(ctx context.Context, tokenStr string) error {
	token, err := u.userRepo.GetVerificationToken(ctx, tokenStr, "EMAIL_VERIFICATION")
	if err != nil {
		return err
	}

	if token.ExpiresAt.Before(time.Now()) {
		return errors.New("verification token expired")
	}

	user, err := u.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return err
	}

	user.IsVerified = true
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	_ = u.userRepo.DeleteVerificationToken(ctx, token.ID)
	return nil
}
