package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
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

	newUser := new(entity.User)
	total, err := c.UserRepository.UserCountByEmail(c.DB, newUser, request.Email)
	if err != nil {
		c.Log.Warnf("failed to count users from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to count users from database : %+v", err))
	}

	if total > 0 {
		c.Log.Warnf("email user has already exists!")
		return nil, fiber.NewError(fiber.StatusConflict, "email user has already exists!")
	}

	if errs := helper_others.ValidatePassword(request.Password, request.Lang); len(errs) > 0 {
		c.Log.Warnf(errs)
		return nil, fiber.NewError(fiber.StatusBadRequest, errs)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("failed to generate bcrypt on password hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to generate bcrypt on password hash : %+v", err))
	}

	newUser.Name.FirstName = request.FirstName
	newUser.Name.LastName = request.LastName
	newUser.Email = request.Email
	newUser.Phone = request.Phone
	newUser.Password = string(password)
	newUser.Role = request.Role
	newUser.VerificationToken = uuid.New().String()
	newUser.TokenExpiry = time.Now().Add(time.Minute * 30)
	newUser.EmailVerified = false
	if err := c.UserRepository.Create(tx, newUser); err != nil {
		c.Log.Warnf("failed to create user into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create user into database : %+v", err))
	}

	// setelah itu buat wallet
	newWallet := &entity.Wallet{}
	newWallet.UserId = newUser.ID
	newWallet.Balance = 0
	newWallet.Status = enum_state.ACTIVE_WALLET
	if err := c.WalletRepository.Create(tx, newWallet); err != nil {
		c.Log.Warnf("failed to create a new wallet into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create a new wallet into database : %+v", err))
	}

	// setelah itu buat cart
	newCart := &entity.Cart{}
	newCart.UserID = newUser.ID
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
	logoImageBase64, err := helper_others.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}

	if request.Role == enum_state.CUSTOMER {
		// setelah semuanya berhasil maka kirim notifikasi email
		newMail := new(model.Mail)
		newMail.To = []string{newUser.Email}
		newMail.Subject = "Email Verification"
		if request.Lang == enum_state.INDONESIA {
			newMail.Subject = "Verifikasi Email"
		}
		baseTemplatePath := "../internal/templates/base_template_email1.html"
		childPath := fmt.Sprintf("../internal/templates/%s/email/email_verification.html", request.Lang)
		tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
		if err != nil {
			c.Log.Warnf("failed to parse template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
		}

		baseURL := fmt.Sprintf("%s://%s/api/users/verify-email/%s", ctx.Protocol(), ctx.Hostname(), newUser.VerificationToken)
		params := url.Values{}
		params.Set("lang", string(request.Lang))
		params.Set("timezone", request.TimeZone.String())

		verifyURL := baseURL
		if encoded := params.Encode(); encoded != "" {
			verifyURL += "?" + encoded
		}

		bodyBuilder := new(strings.Builder)
		err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
			"Name":            newUser.Name.FirstName,
			"Year":            time.Now().Format("2006"),
			"CompanyName":     newApp.AppName,
			"LogoImage":       logoImageBase64,
			"VerificationURL": verifyURL,
			"TokenExpiry":     newUser.TokenExpiry.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
		})
		if err != nil {
			c.Log.Warnf("failed to execute template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
		}

		newMail.Template = *bodyBuilder
		c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
		select {
		case c.Email.MailQueue <- *newMail:
		default:
			c.Log.Warnf("email queue full, failed to send to %s", newUser.Email)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", newUser.Email))
		}
	} else {
		newMail := new(model.Mail)
		newMail.To = []string{newUser.Email}
		newMail.Subject = "Admin Email Verification"
		if request.Lang == enum_state.INDONESIA {
			newMail.Subject = "Verifikasi Email Admin"
		}
		baseTemplatePath := "../internal/templates/base_template_email1.html"
		childPath := fmt.Sprintf("../internal/templates/%s/email/email_verification_admin.html", request.Lang)
		tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
		if err != nil {
			c.Log.Warnf("failed to parse template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
		}

		baseURL := fmt.Sprintf("%s://%s/api/users/verify-email/%s", ctx.Protocol(), ctx.Hostname(), newUser.VerificationToken)
		params := url.Values{}
		params.Set("lang", string(request.Lang))
		params.Set("timezone", request.TimeZone.String())

		verifyURL := baseURL
		if encoded := params.Encode(); encoded != "" {
			verifyURL += "?" + encoded
		}

		bodyBuilder := new(strings.Builder)
		err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
			"AdminName":       newUser.Name.FirstName,
			"LogoImage":       logoImageBase64,
			"CompanyName":     newApp.AppName,
			"VerificationURL": verifyURL,
			"TokenExpiry":     newUser.TokenExpiry.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
		})

		if err != nil {
			c.Log.Warnf("failed to execute template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
		}

		newMail.Template = *bodyBuilder
		c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
		select {
		case c.Email.MailQueue <- *newMail:
		default:
			c.Log.Warnf("email queue full, failed to send to %s", newUser.Email)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", newUser.Email))
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.UserToResponse(newUser), nil
}

func (c *UserUseCase) VerifyEmailRegistration(ctx *fiber.Ctx, request *model.VerifyEmailRegisterRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newUser := new(entity.User)
	if err := c.UserRepository.FindVerifyToken(c.DB, newUser, request.VerificationToken); err != nil {
		c.Log.Warnf("verify token not found : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("verify token not found : %+v", err))
	}

	// cek apakah sudah expired
	if newUser.TokenExpiry.Before(time.Now()) {
		c.Log.Warnf("verify token was expired!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "verify token was expired!")
	}

	if !newUser.EmailVerified {
		newUser.EmailVerified = true
		if err := c.UserRepository.Update(tx, newUser); err != nil {
			c.Log.Warnf("failed to update email verified into database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update email verified into database : %+v", err))
		}

		// send email bahwa registrasi berhasil karena email sudah diverifikasi
		newMail := new(model.Mail)
		newMail.To = []string{newUser.Email}
		newApp := new(entity.Application)
		if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
			c.Log.Warnf("failed to find application from database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
		}

		logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
		logoImageBase64, err := helper_others.ImageToBase64(logoImagePath)
		if err != nil {
			c.Log.Warnf("failed to convert logo to base64 : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
		}
		baseTemplatePath := "../internal/templates/base_template_email1.html"
		childPath := fmt.Sprintf("../internal/templates/%s/email/registration_success.html", request.Lang)
		newMail.Subject = "Registration Successful"
		if request.Lang == enum_state.INDONESIA {
			newMail.Subject = "Registrasi Berhasil"
		}

		if newUser.Role == enum_state.ADMIN {
			childPath = fmt.Sprintf("../internal/templates/%s/email/admin_creation.html", request.Lang)
			newMail.Subject = "New Admin User Created"
			if request.Lang == enum_state.INDONESIA {
				newMail.Subject = "Akun Admin Baru Berhasil Dibuat"
			}
		}

		tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
		if err != nil {
			c.Log.Warnf("failed to parse template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
		}

		loginURL := fmt.Sprintf("%s/auth/login", request.BaseFrontEndURL)
		bodyBuilder := new(strings.Builder)
		err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
			"FirstName":   newUser.Name.FirstName,
			"LastName":    newUser.Name.LastName,
			"LoginURL":    loginURL,
			"Year":        time.Now().Format("2006"),
			"CompanyName": newApp.AppName,
			"LogoImage":   logoImageBase64,
			"Email":       newUser.Email,
			"CreatedAt":   newUser.CreatedAt.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
		})
		if err != nil {
			c.Log.Warnf("failed to execute template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
		}

		newMail.Template = *bodyBuilder
		c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
		select {
		case c.Email.MailQueue <- *newMail:
		default:
			c.Log.Warnf("email queue full, failed to send to %s", newUser.Email)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", newUser.Email))
		}

		// send notifikasi bahwa registrasi berhasil karena email sudah diverifikasi
		newNotification := new(entity.Notification)
		newNotification.Title = newMail.Subject
		newNotification.UserID = newUser.ID
		newNotification.IsRead = false
		newNotification.Type = enum_state.AUTHENTICATION
		baseTemplatePath = "../internal/templates/base_template_notification1.html"
		childPath = fmt.Sprintf("../internal/templates/%s/notification/registration_success.html", request.Lang)
		tmpl, err = template.ParseFiles(baseTemplatePath, childPath)
		if err != nil {
			c.Log.Warnf("failed to parse template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
		}

		logoImagePath = fmt.Sprintf("%s://%s/api/image/application/%s", ctx.Protocol(), ctx.Hostname(), newApp.LogoFilename)
		bodyBuilder = new(strings.Builder)
		err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
			"FirstName":     newUser.Name.FirstName,
			"Year":          time.Now().Format("2006"),
			"CompanyName":   newApp.AppName,
			"LogoImagePath": logoImagePath,
		})

		if err != nil {
			c.Log.Warnf("failed to execute template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
		}
		newNotification.BodyContent = bodyBuilder.String()

		if err := c.NotificationRepository.Create(tx, newNotification); err != nil {
			c.Log.Warnf("failed to create notification into database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create notification into database : %+v", err))
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.UserToResponse(newUser), nil
}

func (c *UserUseCase) ValidateVerifyTokenIsValid(ctx *fiber.Ctx, verifyToken string, email string) error {
	tx := c.DB.WithContext(ctx.Context())

	newUser := new(entity.User)
	if err := c.UserRepository.FindVerifyToken(tx, newUser, verifyToken); err != nil {
		c.Log.Warnf("verify token not found : %+v", err)
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("verify token not found : %+v", err))
	}

	if !newUser.EmailVerified {
		c.Log.Warnf("email not verified!")
		return fiber.NewError(fiber.StatusBadRequest, "email not verified!")
	}

	if newUser.Email != email {
		c.Log.Warnf("email is not match!")
		return fiber.NewError(fiber.StatusBadRequest, "email is not match!")
	}

	return nil
}

func (c *UserUseCase) Authenticate(ctx context.Context, request *model.LoginUserRequest) (*model.UserTokenResponse, *model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newUser := new(entity.User)
	if err := c.UserRepository.FindByEmail(c.DB, newUser, request.Email); err != nil {
		c.Log.Warnf("user not found : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("user not found : %+v", err))
	}

	if !newUser.EmailVerified {
		c.Log.Warnf("your account has not verified email")
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, "your account has not verified email!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("password is wrong : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("password is wrong : %+v", err))
	}

	var token = &entity.Token{}
	now := time.Now()
	hours := 24
	if request.Remember {
		hours = hours * 3
	}
	totalHours := now.Add(time.Duration(hours) * time.Hour)
	token.Token = uuid.New().String()
	token.UserId = newUser.ID
	token.ExpiryDate = totalHours
	if err := c.TokenRepository.Create(tx, token); err != nil {
		c.Log.Warnf("failed to create token by user into database : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create token by user into database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.UserTokenToResponse(token), converter.UserToResponse(newUser), nil
}

func (c *UserUseCase) GetUserByToken(ctx context.Context, request *model.GetUserByTokenRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("token is not included in header : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("token is not included in header : %+v", err))
	}

	user := new(entity.User)
	if err := c.UserRepository.FindUserByToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("token isn't valid!")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "token isn't valid!")
	}

	expiredDate := user.Token.ExpiryDate
	fmt.Println("current time:", time.Now())
	fmt.Println("expired date:", expiredDate)
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
	user.Role = enum_state.ADMIN

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

	if newUser.ID == 0 {
		c.Log.Warnf("user not found!")
		return false, fiber.NewError(fiber.StatusUnauthorized, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(request.OldPassword)); err != nil {
		c.Log.Warnf("old Password is wrong : %+v", err)
		return false, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("old Password is wrong : %+v", err))
	}

	// lalu hapus usernya (soft delete)
	if err := c.UserRepository.Delete(tx, newUser); err != nil {
		c.Log.Warnf("can't delete current user : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't delete current user : %+v", err))
	}

	// kirim email
	newApp := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
		c.Log.Warnf("failed to find application from database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
	}

	logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
	logoImageBase64, err := helper_others.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}

	newMail := new(model.Mail)
	newMail.To = []string{newUser.Email}
	baseTemplatePath := "../internal/templates/base_template_email1.html"
	childPath := fmt.Sprintf("../internal/templates/%s/email/account_deletion.html", request.Lang)
	newMail.Subject = "Account Deletion Successful"
	if request.Lang == enum_state.INDONESIA {
		newMail.Subject = "Penghapusan Akun Berhasil"
	}
	tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
		"FirstName":   newUser.Name.FirstName,
		"Year":        time.Now().Format("2006"),
		"CompanyName": newApp.AppName,
		"LogoImage":   logoImageBase64,
		"DeletedAt":   newUser.DeletedAt.Time.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
	})
	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newMail.Template = *bodyBuilder
	c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
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

	totalMinute := time.Minute * 5
	code := rand.Intn(900000) + 100000
	newPasswordReset.VerificationCode = code
	newPasswordReset.ExpiresAt = time.Now().Add(totalMinute)
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
	logoImageBase64, err := helper_others.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}
	baseTemplatePath := "../internal/templates/base_template_email1.html"
	childPath := fmt.Sprintf("../internal/templates/%s/email/forgot_password.html", request.Lang)
	newMail.Subject = "Forgot Password"
	if request.Lang == enum_state.INDONESIA {
		newMail.Subject = "Lupa Kata Sandi"
	}

	tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]any{
		"FirstName":          newUser.Name.FirstName,
		"Year":               time.Now().Format("2006"),
		"CompanyName":        newApp.AppName,
		"LogoImage":          logoImageBase64,
		"VerificationCode":   strconv.Itoa(code),
		"TotalMinuteExpired": totalMinute.Minutes(),
	})
	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newMail.Template = *bodyBuilder
	c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
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
		c.Log.Warnf("password reset or verification code not found!")
		return false, fiber.NewError(fiber.StatusNotFound, "password reset or verification code not found!")
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

	if errs := helper_others.ValidatePassword(request.NewPassword, request.Lang); len(errs) > 0 {
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

	logoImagePathURL := fmt.Sprintf("%s://%s/api/image/application/%s", ctx.Protocol(), ctx.Hostname(), newApp.LogoFilename)
	baseTemplatePath := "../internal/templates/base_template_notification1.html"
	childPath := fmt.Sprintf("../internal/templates/%s/notification/password_reset.html", request.Lang)
	tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
		"FirstName":     newUser.Name.FirstName,
		"Year":          time.Now().Format("2006"),
		"CompanyName":   newApp.AppName,
		"LogoImagePath": logoImagePathURL,
		"UpdatedAt":     newUser.UpdatedAt.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
	})

	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newNotification := new(entity.Notification)
	newNotification.UserID = newUser.ID
	newNotification.Title = "Password Reset Successful"
	if request.Lang == enum_state.INDONESIA {
		newNotification.Title = "Mengatur Ulang Kata Sandi Berhasil"
	}
	newNotification.IsRead = false
	newNotification.Type = enum_state.AUTHENTICATION
	newNotification.BodyContent = bodyBuilder.String()
	if err := c.NotificationRepository.Create(tx, newNotification); err != nil {
		c.Log.Warnf("failed to create notification into database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create notification into database : %+v", err))
	}

	// send email
	newMail := new(model.Mail)
	newMail.To = []string{newUser.Email}
	newMail.Cc = []string{}
	baseTemplatePath = "../internal/templates/base_template_email1.html"
	childPath = fmt.Sprintf("../internal/templates/%s/email/password_reset.html", request.Lang)
	newMail.Subject = "Password Reset Successful"
	if request.Lang == enum_state.INDONESIA {
		newMail.Subject = "Atur Ulang Kata Sandi Berhasil"
	}
	tmpl, err = template.ParseFiles(baseTemplatePath, childPath)
	if err != nil {
		c.Log.Warnf("failed to parse template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
	}

	logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
	logoImageBase64, err := helper_others.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}

	bodyBuilder = new(strings.Builder)
	err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
		"FirstName":   newUser.Name.FirstName,
		"Year":        time.Now().Format("2006"),
		"CompanyName": newApp.AppName,
		"LogoImage":   logoImageBase64,
		"DeletedAt":   newUser.UpdatedAt.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
	})
	if err != nil {
		c.Log.Warnf("failed to execute template file html : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
	}

	newMail.Template = *bodyBuilder
	c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
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
