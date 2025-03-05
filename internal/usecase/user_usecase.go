package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	UserRepository    *repository.UserRepository
	TokenRepository   *repository.TokenRepository
	AddressRepository *repository.AddressRepository
	WalletRepository  *repository.WalletRepository
	CartRepository    *repository.CartRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, tokenRepository *repository.TokenRepository,
	addressRepository *repository.AddressRepository, walletRepository *repository.WalletRepository,
	cartRepository *repository.CartRepository) *UserUseCase {
	return &UserUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		UserRepository:    userRepository,
		TokenRepository:   tokenRepository,
		AddressRepository: addressRepository,
		WalletRepository:  walletRepository,
		CartRepository:    cartRepository,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %v", err))
	}

	user := &entity.User{}
	total, err := c.UserRepository.UserCountByEmail(c.DB, user, request.Email)
	if err != nil {
		c.Log.Warnf("Failed count users from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed count users from database : %+v", err))
	}

	if total > 0 {
		c.Log.Warnf("Email user has already exists!",)
		return nil, fiber.NewError(fiber.StatusConflict, "Email user has already exists!")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt on password hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to generate bcrypt on password hash : %+v", err))
	}

	user.Name.FirstName = request.FirstName
	user.Name.LastName = request.LastName
	user.Email = request.Email
	user.Phone = request.Phone
	user.Password = string(password)
	user.Role = request.Role
	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed create user into database : %+v", err))
	}

	// setelah itu buat wallet
	newWallet := &entity.Wallet{}
	newWallet.UserId = user.ID
	newWallet.Balance = 0
	newWallet.Status = helper.ACTIVE
	if err := c.WalletRepository.Create(tx, newWallet); err != nil {
		c.Log.Warnf("Failed to create a new wallet into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create a new wallet into database : %+v", err))
	}

	// setelah itu buat cart
	newCart := &entity.Cart{}
	newCart.UserID = user.ID
	if err := c.CartRepository.Create(tx, newCart); err != nil {
		c.Log.Warnf("Failed to create a new cart into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create a new cart into database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequst) (*model.UserTokenResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(c.DB, user, request.Email); err != nil {
		c.Log.Warnf("User not found : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("User not found : %+v", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Password is wrong : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Password is wrong : %+v", err))
	}

	var token = &entity.Token{}
	now := time.Now()
	oneHours := now.Add(24 * time.Hour)
	token.Token = uuid.New().String()
	token.UserId = user.ID
	token.ExpiryDate = oneHours
	if err := c.TokenRepository.Create(tx, token); err != nil {
		c.Log.Warnf("Failed to create token by user into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create token by user into database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed commit transaction : %+v", err))
	}

	return converter.UserTokenToResponse(token), nil
}

func (c *UserUseCase) GetUserByToken(ctx context.Context, request *model.GetUserByTokenRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Token is not included in header : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindUserByToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("Token isn't valid : %+v", err)
		return nil, fiber.ErrUnauthorized
	}
	
	expiredDate := user.Token.ExpiryDate
	if expiredDate.Before(time.Now()) {
		c.Log.Warn("Token is expired")
		return nil, fiber.ErrUnauthorized
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest, currentUser *model.UserResponse) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	totalCount, err := c.UserRepository.CheckEmailIsExists(tx, currentUser.Email, request.Email)
	if err != nil {
		c.Log.Warnf("Cannot count email is exists : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if totalCount > 0 {
		c.Log.Warnf("Email has already exists : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	user.ID = currentUser.ID
	user.Email = request.Email
	user.Name.FirstName = request.FirstName
	user.Name.LastName = request.LastName
	user.Phone = request.Phone
	user.Role = helper.ADMIN

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed to update data user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) UpdatePassword(ctx context.Context, request *model.UpdateUserPasswordRequest, user *model.UserResponse) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newUser := new(entity.User)
	if err := c.UserRepository.FindUserById(tx, newUser, user.ID); err != nil {
		c.Log.Warnf("Token isn't valid : %+v", err)
		return false, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(request.OldPassword)); err != nil {
		c.Log.Warnf("Old Password is wrong : %+v", err)
		return false, fiber.ErrUnauthorized
	}

	newPasswordRequest, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return false, fiber.ErrInternalServerError
	}
	newUser.Password = string(newPasswordRequest)

	if err := c.UserRepository.Update(tx, newUser); err != nil {
		c.Log.Warnf("Failed to update data user : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}
	return true, nil
}

func (c *UserUseCase) Logout(ctx context.Context, token *model.GetUserByTokenRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newToken := new(entity.Token)
	result := c.TokenRepository.DeleteToken(tx, newToken, token.Token)
	if result.RowsAffected == 0 {
		c.Log.Warnf("Can't delete token : %+v", result.Error)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func (c *UserUseCase) RemoveCurrentAccount(ctx context.Context, request *model.DeleteCurrentUserRequest, user *model.UserResponse) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newUser := new(entity.User)
	if err := c.UserRepository.FindUserById(tx, newUser, user.ID); err != nil {
		c.Log.Warnf("Can't find user by token : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(request.OldPassword)); err != nil {
		c.Log.Warnf("Old Password is wrong : %+v", err)
		return false, fiber.ErrUnauthorized
	}

	// hapus token terlebih dahulu
	newToken := new(entity.Token)
	deleteToken := c.TokenRepository.DeleteToken(tx, newToken, newUser.Token.Token)
	if deleteToken.RowsAffected == 0 {
		c.Log.Warnf("Can't delete token : %+v", deleteToken.Error)
		return false, fiber.ErrInternalServerError
	}

	// hapus address terlebih dahulu
	newAddress := new(entity.Address)
	if err := c.AddressRepository.DeleteAllAddressByUserId(tx, newAddress, newUser.ID); err.Error != nil {
		c.Log.Warnf("Can't delete addresses by user id : %+v", err.Error)
		return false, fiber.ErrInternalServerError
	}

	// lalu hapus usernya
	if err := c.UserRepository.Delete(tx, newUser); err != nil {
		c.Log.Warnf("Can't delete current user : %+v", deleteToken.Error)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
