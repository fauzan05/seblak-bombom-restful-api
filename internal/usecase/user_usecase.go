package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/helper/mailer"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	UserRepository         *repository.UserRepository
	TokenRepository        *repository.TokenRepository
	AddressRepository      *repository.AddressRepository
	WalletRepository       *repository.WalletRepository
	CartRepository         *repository.CartRepository
	NotificationRepository *repository.NotificationRepository
	Email                  *mailer.EmailWorker
	ApplicationRepository  *repository.ApplicationRepository
	PasswordReset          *repository.PasswordResetRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, tokenRepository *repository.TokenRepository,
	addressRepository *repository.AddressRepository, walletRepository *repository.WalletRepository,
	cartRepository *repository.CartRepository, notificationRepository *repository.NotificationRepository,
	email *mailer.EmailWorker, applicationRepository *repository.ApplicationRepository,
	passwordReset *repository.PasswordResetRepository) *UserUseCase {
	return &UserUseCase{
		DB:                     db,
		Log:                    log,
		Validate:               validate,
		UserRepository:         userRepository,
		TokenRepository:        tokenRepository,
		AddressRepository:      addressRepository,
		WalletRepository:       walletRepository,
		CartRepository:         cartRepository,
		NotificationRepository: notificationRepository,
		Email:                  email,
		ApplicationRepository:  applicationRepository,
		PasswordReset:          passwordReset,
	}
}

func (c *UserUseCase) Create(ctx *fiber.Ctx, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %v", err))
	}

	user := new(entity.User)
	total, err := c.UserRepository.UserCountByEmail(c.DB, user, request.Email)
	if err != nil {
		c.Log.Warnf("failed to count users from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to count users from database : %+v", err))
	}

	if total > 0 {
		c.Log.Warnf("email user has already exists!")
		return nil, fiber.NewError(fiber.StatusConflict, "email user has already exists!")
	}

	if errs := helper.ValidatePassword(request.Password); len(errs) > 0 {
		c.Log.Warnf(errs)
		return nil, fiber.NewError(fiber.StatusBadRequest, errs)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("failed to generate bcrypt on password hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to generate bcrypt on password hash : %+v", err))
	}

	user.Name.FirstName = request.FirstName
	user.Name.LastName = request.LastName
	user.Email = request.Email
	user.Phone = request.Phone
	user.Password = string(password)
	user.Role = request.Role
	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("failed to create user into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create user into database : %+v", err))
	}

	// setelah itu buat wallet
	newWallet := &entity.Wallet{}
	newWallet.UserId = user.ID
	newWallet.Balance = 0
	newWallet.Status = helper.ACTIVE_WALLET
	if err := c.WalletRepository.Create(tx, newWallet); err != nil {
		c.Log.Warnf("failed to create a new wallet into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create a new wallet into database : %+v", err))
	}

	// setelah itu buat cart
	newCart := &entity.Cart{}
	newCart.UserID = user.ID
	if err := c.CartRepository.Create(tx, newCart); err != nil {
		c.Log.Warnf("failed to create a new cart into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create a new cart into database : %+v", err))
	}

	newApp := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
		c.Log.Warnf("failed to find application from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
	}

	logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
	logoImageBase64, err := helper.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}
	newNotification := new(entity.Notification)
	newNotification.UserID = user.ID
	newNotification.Title = "Registration Successful 🎉"
	newNotification.Message = fmt.Sprintf("Hi %s, your account is now active. Welcome!", user.Name.FirstName)
	templatePath := "../internal/templates/english/notification/registration_success.html"
	if request.Language == helper.INDONESIA {
		newNotification.Title = "Registrasi Berhasil 🎉"
		newNotification.Message = fmt.Sprintf("Hai %s, akunmu sekarang sudah aktif. Selamat Datang!", user.Name.FirstName)
		templatePath = "../internal/templates/indonesia/notification/registration_success.html"
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]string{
		"Name":        user.Name.FirstName,
		"Year":        time.Now().Format("2006"),
		"CompanyName": newApp.AppName,
		"LogoImage":   logoImageBase64,
	})
	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newNotification.IsRead = false
	newNotification.Type = helper.AUTHENTICATION
	newNotification.BodyContent = bodyBuilder.String()
	if err := c.NotificationRepository.Create(tx, newNotification); err != nil {
		c.Log.Warnf("failed to create notification into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create notification into database : %+v", err))
	}

	// setelah semuanya berhasil maka kirim notifikasi email
	newMail := new(model.Mail)
	newMail.To = []string{user.Email}
	newMail.Cc = []string{}
	newMail.Subject = "Registration Successful"
	templatePath = "../internal/templates/english/email/registration_success.html"
	if request.Language == helper.INDONESIA {
		newMail.Subject = "Registrasi Berhasil"
		templatePath = "../internal/templates/indonesia/email/registration_success.html"
	}
	tmpl, err = template.ParseFiles(templatePath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder = new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]string{
		"Name":        user.Name.FirstName,
		"LoginURL":    "http://localhost:8000/login",
		"Year":        time.Now().Format("2006"),
		"CompanyName": newApp.AppName,
		"LogoImage":   logoImageBase64,
	})
	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newMail.Template = *bodyBuilder
	select {
	case c.Email.MailQueue <- *newMail:
	default:
		c.Log.Warnf("email queue full, failed to send to %s", user.Email)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", user.Email))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Authenticate(ctx context.Context, request *model.LoginUserRequest) (*model.UserTokenResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(c.DB, user, request.Email); err != nil {
		c.Log.Warnf("user not found : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("user not found : %+v", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("password is wrong : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("password is wrong : %+v", err))
	}

	var token = &entity.Token{}
	now := time.Now()
	oneHours := now.Add(24 * time.Hour)
	token.Token = uuid.New().String()
	token.UserId = user.ID
	token.ExpiryDate = oneHours
	if err := c.TokenRepository.Create(tx, token); err != nil {
		c.Log.Warnf("failed to create token by user into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create token by user into database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.UserTokenToResponse(token), nil
}

func (c *UserUseCase) GetUserByToken(ctx context.Context, request *model.GetUserByTokenRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("token is not included in header : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("token is not included in header : %+v", err))
	}

	user := new(entity.User)
	if err := c.UserRepository.FindUserByToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("token isn't valid : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("token isn't valid : %+v", err))
	}

	expiredDate := user.Token.ExpiryDate
	if expiredDate.Before(time.Now()) {
		c.Log.Warn("token is expired!")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "token is expired!")
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest, currentUser *model.UserResponse) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	totalCount, err := c.UserRepository.CheckEmailIsExists(tx, currentUser.Email, request.Email)
	if err != nil {
		c.Log.Warnf("failed to count email is exist : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to count email is exist : %+v", err))
	}

	if totalCount > 0 {
		c.Log.Warnf("email has already exists!")
		return nil, fiber.NewError(fiber.StatusConflict, "email has already exists!")
	}

	user := new(entity.User)
	user.ID = currentUser.ID
	user.Email = request.Email
	user.Name.FirstName = request.FirstName
	user.Name.LastName = request.LastName
	user.Phone = request.Phone
	user.Role = helper.ADMIN

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("failed to update data user : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update data user : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) UpdatePassword(ctx context.Context, request *model.UpdateUserPasswordRequest, user *model.UserResponse) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newUser := new(entity.User)
	if err := c.UserRepository.FindUserById(tx, newUser, user.ID); err != nil {
		c.Log.Warnf("token isn't valid : %+v", err)
		return false, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("token isn't valid : %+v", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(request.OldPassword)); err != nil {
		c.Log.Warnf("old Password is wrong : %+v", err)
		return false, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("old Password is wrong : %+v", err))
	}

	newPasswordRequest, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("failed to generate bcrypt hash : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to generate bcrypt hash : %+v", err))
	}

	newUser.Password = string(newPasswordRequest)
	if err := c.UserRepository.Update(tx, newUser); err != nil {
		c.Log.Warnf("failed to update password user : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update password user : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}

func (c *UserUseCase) Logout(ctx context.Context, token *model.GetUserByTokenRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newToken := new(entity.Token)
	result := c.TokenRepository.DeleteToken(tx, newToken, token.Token)
	if result.RowsAffected == 0 {
		c.Log.Warnf("can't delete token : %+v", result.Error)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't delete token : %+v", result.Error))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}

func (c *UserUseCase) RemoveCurrentAccount(ctx context.Context, request *model.DeleteCurrentUserRequest, user *model.UserResponse) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newUser := new(entity.User)
	if err := c.UserRepository.FindUserById(tx, newUser, user.ID); err != nil {
		c.Log.Warnf("can't find user by token : %+v", err)
		return false, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("can't find user by token : %+v", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(request.OldPassword)); err != nil {
		c.Log.Warnf("old Password is wrong : %+v", err)
		return false, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("old Password is wrong : %+v", err))
	}

	// hapus token terlebih dahulu
	newToken := new(entity.Token)
	deleteToken := c.TokenRepository.DeleteToken(tx, newToken, newUser.Token.Token)
	if deleteToken.RowsAffected == 0 {
		c.Log.Warnf("can't delete token : %+v", deleteToken.Error)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't delete token : %+v", deleteToken.Error))
	}

	// hapus address terlebih dahulu
	newAddress := new(entity.Address)
	if err := c.AddressRepository.DeleteAllAddressByUserId(tx, newAddress, newUser.ID); err.Error != nil {
		c.Log.Warnf("can't delete addresses by user id : %+v", err.Error)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't delete addresses by user id : %+v", err.Error))
	}

	// lalu hapus usernya
	if err := c.UserRepository.Delete(tx, newUser); err != nil {
		c.Log.Warnf("can't delete current user : %+v", deleteToken.Error)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't delete current user : %+v", deleteToken.Error))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}

func (c *UserUseCase) AddForgotPassword(ctx *fiber.Ctx, request *model.CreateForgotPassword) (*model.PasswordResetResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newUser := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, newUser, request.Email); err != nil {
		c.Log.Warnf("failed to find email address : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to find email address : %+v", err))
	}

	newPasswordReset := new(entity.PasswordReset)
	newPasswordReset.UserId = newUser.ID
	count, err := c.PasswordReset.FindAndCountEntityByUserId(tx, newPasswordReset, newUser.ID)
	if err != nil {
		c.Log.Warnf("failed to find password reset from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find password reset from database : %+v", err))
	}

	if count > 0 {
		// maka delete dulu
		if err := c.PasswordReset.Delete(tx, newPasswordReset); err != nil {
			c.Log.Warnf("failed to delete password reset into database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete password reset into database : %+v", err))
		}
	}

	code := rand.Intn(900000) + 100000
	newPasswordReset.VerificationCode = code
	newPasswordReset.ExpiresAt = time.Now().Add(time.Minute * 5)
	if err := c.PasswordReset.Create(tx, newPasswordReset); err != nil {
		c.Log.Warnf("failed to create category into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed create category into database : %+v", err))
	}

	// kirim email
	newMail := new(model.Mail)
	newMail.To = []string{newUser.Email}
	newMail.Cc = []string{}
	newApp := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
		c.Log.Warnf("failed to find application from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
	}

	logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
	logoImageBase64, err := helper.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}
	newMail.Subject = "Forgot Password"
	templatePath := "../internal/templates/english/email/forgot_password.html"
	if request.Lang == helper.INDONESIA {
		newMail.Subject = "Lupa Kata Sandi"
		templatePath = "../internal/templates/indonesia/email/forgot_password.html"
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]string{
		"Name":               newUser.Name.FirstName,
		"Year":               time.Now().Format("2006"),
		"CompanyName":        newApp.AppName,
		"LogoImage":          logoImageBase64,
		"VerificationCode":   strconv.Itoa(code),
		"TotalMinuteExpired": "5",
	})
	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newMail.Template = *bodyBuilder
	select {
	case c.Email.MailQueue <- *newMail:
	default:
		c.Log.Warnf("email queue full, failed to send to %s", newUser.Email)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", newUser.Email))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.PasswordResetToResponse(newPasswordReset), nil
}

func (c *UserUseCase) ValidateForgotPassword(ctx *fiber.Ctx, request *model.ValidateForgotPassword) (bool, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newPasswordReset := new(entity.PasswordReset)
	newPasswordReset.ID = request.ID
	count, err := c.PasswordReset.FindAndCountById(tx, newPasswordReset)
	if err != nil {
		c.Log.Warnf("failed to find password reset from database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find password reset from database : %+v", err))
	}

	if count < 1 {
		c.Log.Warnf("password reset not found!")
		return false, fiber.NewError(fiber.StatusNotFound, "password reset not found!")
	}

	// cek apakah valid dan tanggal belum expired
	if newPasswordReset.ExpiresAt.Before(time.Now()) {
		// maka return error bahwa verification code expired
		c.Log.Warnf("password reset was expired!")
		return false, fiber.NewError(fiber.StatusBadRequest, "password reset was expired!")
	}

	if newPasswordReset.VerificationCode != request.VerificationCode {
		c.Log.Warnf("verification code is not match!")
		return false, fiber.NewError(fiber.StatusBadRequest, "verification code is not match!")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}

func (c *UserUseCase) Reset(ctx *fiber.Ctx, request *model.PasswordResetRequest) (bool, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newPasswordReset := new(entity.PasswordReset)
	newPasswordReset.ID = request.ID
	count, err := c.PasswordReset.FindAndCountById(tx, newPasswordReset)
	if err != nil {
		c.Log.Warnf("failed to find password reset from database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find password reset from database : %+v", err))
	}

	if count < 1 {
		c.Log.Warnf("password reset not found!")
		return false, fiber.NewError(fiber.StatusNotFound, "password reset not found!")
	}

	// cek apakah valid dan tanggal belum expired
	if newPasswordReset.ExpiresAt.Add(time.Minute * 15).Before(time.Now()) {
		// maka return error bahwa verification code expired
		c.Log.Warnf("password reset was expired!")
		return false, fiber.NewError(fiber.StatusBadRequest, "password reset was expired!")
	}

	if newPasswordReset.VerificationCode != request.VerificationCode {
		c.Log.Warnf("verification code is not match!")
		return false, fiber.NewError(fiber.StatusBadRequest, "verification code is not match!")
	}

	if errs := helper.ValidatePassword(request.NewPassword); len(errs) > 0 {
		c.Log.Warnf(errs)
		return false, fiber.NewError(fiber.StatusBadRequest, errs)
	}

	newUser := new(entity.User)
	newUser.ID = newPasswordReset.UserId
	if err := c.UserRepository.FindFirst(tx, newUser); err != nil {
		c.Log.Warnf("failed to find user : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to find user : %+v", err))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("failed to generate bcrypt on password hash : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to generate bcrypt on password hash : %+v", err))
	}
	newUser.Password = string(password)
	if err := c.UserRepository.Update(tx, newUser); err != nil {
		c.Log.Warnf("failed to update password user : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update password user : %+v", err))
	}

	// send notif
	newApp := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
		c.Log.Warnf("failed to find application from database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
	}

	// logoImagePath := fmt.Sprintf("%s://%s/api/image/application/%s", ctx.Protocol(), ctx.Hostname(), newApp.LogoFilename)
	logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
	logoImageBase64, err := helper.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}
	templatePath := "../internal/templates/english/notification/password_reset.html"
	if request.Lang == helper.INDONESIA {
		templatePath = "../internal/templates/indonesia/notification/password_reset.html"
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]string{
		"Name":        newUser.Name.FirstName,
		"Year":        time.Now().Format("2006"),
		"CompanyName": newApp.AppName,
		"LogoImage":   logoImageBase64,
	})

	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newNotification := new(entity.Notification)
	newNotification.UserID = newUser.ID
	newNotification.Title = "Password Reset Successful 🎉"
	newNotification.Message = fmt.Sprintf("Hi %s, We've successfully updated your password. You can now log in with your new credentials. If you did not request this change, please contact our support immediately to secure your account.", newUser.Name.FirstName)
	if request.Lang == helper.INDONESIA {
		newNotification.Title = "Mengatur Ulang Kata Sandi Berhasil 🎉"
		newNotification.Message = fmt.Sprintf("Hai %s, Kata sandi Anda berhasil diperbarui. Sekarang Anda dapat masuk menggunakan kata sandi baru Anda. Jika Anda tidak meminta perubahan ini, segera hubungi tim dukungan kami untuk mengamankan akun Anda.", newUser.Name.FirstName)
	}
	newNotification.IsRead = false
	newNotification.Type = helper.AUTHENTICATION
	newNotification.BodyContent = bodyBuilder.String()
	if err := c.NotificationRepository.Create(tx, newNotification); err != nil {
		c.Log.Warnf("failed to create notification into database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create notification into database : %+v", err))
	}

	// send email
	newMail := new(model.Mail)
	newMail.To = []string{newUser.Email}
	newMail.Cc = []string{}
	newMail.Subject = "Password Reset Successful"
	templatePath = "../internal/templates/english/email/password_reset.html"
	if request.Lang == helper.INDONESIA {
		newMail.Subject = "Atur Ulang Kata Sandi Berhasil"
		templatePath = "../internal/templates/indonesia/email/password_reset.html"
	}
	tmpl, err = template.ParseFiles(templatePath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder = new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]string{
		"Name":        newUser.Name.FirstName,
		"Year":        time.Now().Format("2006"),
		"CompanyName": newApp.AppName,
		"LogoImage":   logoImageBase64,
	})
	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newMail.Template = *bodyBuilder
	select {
	case c.Email.MailQueue <- *newMail:
	default:
		c.Log.Warnf("email queue full, failed to send to %s", newUser.Email)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", newUser.Email))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}
