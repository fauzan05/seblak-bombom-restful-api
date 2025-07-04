package helper_others

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func PrintStructFields(s any) {
	v := reflect.ValueOf(s)

	// Cek apakah s adalah pointer, jika iya, gunakan Elem()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for i := range v.NumField() {
		field := t.Field(i).Name
		value := v.Field(i).Interface()
		fmt.Printf("%s: %v\n", field, value)
	}
}

func SetFiberStatusCode(statusCode string) int {
	// Mapping HTTP status ke fiber error
	errorMap := map[string]*fiber.Error{
		"400": fiber.ErrBadRequest,            // Bad Request
		"401": fiber.ErrUnauthorized,          // Unauthorized
		"403": fiber.ErrForbidden,             // Forbidden
		"404": fiber.ErrNotFound,              // Not Found
		"405": fiber.ErrMethodNotAllowed,      // Method Not Allowed
		"408": fiber.ErrRequestTimeout,        // Request Timeout
		"409": fiber.ErrConflict,              // Conflict
		"413": fiber.ErrRequestEntityTooLarge, // Payload Too Large
		"415": fiber.ErrUnsupportedMediaType,  // Unsupported Media Type
		"422": fiber.ErrUnprocessableEntity,   // Unprocessable Entity
		"429": fiber.ErrTooManyRequests,       // Too Many Requests
		"500": fiber.ErrInternalServerError,   // Internal Server Error
		"501": fiber.ErrNotImplemented,        // Not Implemented
		"502": fiber.ErrBadGateway,            // Bad Gateway
		"503": fiber.ErrServiceUnavailable,    // Service Unavailable
		"504": fiber.ErrGatewayTimeout,        // Gateway Timeout
	}

	// Ambil error dari map, default ke Internal Server Error jika tidak ada
	fiberErr, ok := errorMap[statusCode]
	if !ok {
		fiberErr = fiber.ErrInternalServerError
	}

	return fiberErr.Code
}

type TimeRFC3339 time.Time

func (t TimeRFC3339) MarshalJSON() ([]byte, error) {
	stamp := time.Time(t).Format(time.RFC3339)
	return []byte(`"` + stamp + `"`), nil
}

func (t *TimeRFC3339) UnmarshalJSON(b []byte) error {
	parsed, err := time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err != nil {
		return err
	}
	*t = TimeRFC3339(parsed)
	return nil
}

func (t TimeRFC3339) ToTime() time.Time {
	return time.Time(t)
}

func RoundFloat32(val float32, precision int) float32 {
	f64 := float64(val)
	multiplier := math.Pow(10, float64(precision))
	rounded := math.Round(f64*multiplier) / multiplier
	return float32(rounded)
}

// GenerateBoundary membuat boundary unik untuk multipart email
func GenerateBoundary() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		// fallback kalau gagal
		return "BOUNDARY_DEFAULT"
	}
	return "BOUNDARY_" + hex.EncodeToString(bytes)
}

func ValidatePassword(password string, lang enum_state.Languange) string {
	var errs string

	if lang == enum_state.INDONESIA {
		if len(password) < 8 {
			errs += "Password harus terdiri dari minimal 8 karakter;"
		}
		if len(password) > 100 {
			errs += "Password tidak boleh lebih dari 100 karakter;"
		}
		if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
			errs += "Password harus mengandung setidaknya satu huruf kapital (A-Z);"
		}
		if !regexp.MustCompile(`[a-z]`).MatchString(password) {
			errs += "Password harus mengandung setidaknya satu huruf kecil (a-z);"
		}
		if !regexp.MustCompile(`[0-9]`).MatchString(password) {
			errs += "Password harus mengandung setidaknya satu angka (0-9);"
		}
		if !regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(password) {
			errs += "Password harus mengandung setidaknya satu simbol (!@#~$%^&*()+|_);"
		}
	} else {
		if len(password) < 8 {
			errs += "Password must be at least 8 characters long;"
		}
		if len(password) > 100 {
			errs += "Password must not exceed 100 characters;"
		}
		if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
			errs += "Password must contain at least one uppercase letter;"
		}
		if !regexp.MustCompile(`[a-z]`).MatchString(password) {
			errs += "Password must contain at least one lowercase letter;"
		}
		if !regexp.MustCompile(`[0-9]`).MatchString(password) {
			errs += "Password must contain at least one number;"
		}
		if !regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(password) {
			errs += "Password must contain at least one symbol (!@#~$%^&*()+|_);"
		}
	}

	return errs
}

func ImageToBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func GetPaymentStatusColor(status enum_state.PaymentStatus) string {
	switch status {
	case enum_state.PAID_PAYMENT:
		return "green"
	case enum_state.PENDING_PAYMENT:
		return "orange"
	case enum_state.CANCELLED_PAYMENT, enum_state.FAILED_PAYMENT, enum_state.EXPIRED_PAYMENT:
		return "red"
	default:
		return "gray" // atau warna lain untuk status yang tidak diketahui
	}
}

var TimeZoneMap = map[string]string{
	"Asia/Jakarta":   "WIB",
	"Asia/Pontianak": "WIB",
	"Asia/Makassar":  "WITA",
	"Asia/Jayapura":  "WIT",
	"UTC":            "UTC",
	// tambahkan time zone lainnya di Indonesia
}

func FormatNumberFloat32(n float32) string {
	str := strconv.FormatFloat(float64(n), 'f', 0, 32)
	parts := strings.Split(str, ".")
	var integerPart string
	var fractionalPart string

	integerPart = parts[0]
	if len(parts) > 1 {
		fractionalPart = parts[1]
	}

	var result string
	for i, r := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			result += "."
		}
		result += string(r)
	}

	if fractionalPart != "" {
		result += "," + fractionalPart
	}

	return result
}

type SaveWalletTransactionRequest struct {
	DB              *gorm.DB
	UserId          uint64
	OrderId         *uint64
	Amount          float32
	FlowType        enum_state.WalletFlowType
	TransactionType enum_state.WalletTransactionType
	PaymentMethod   enum_state.PaymentMethod
	Status          enum_state.WalletTransactionStatus
	ReferenceNumber string
	Note            string
	AdminNote       string
	ProcessedBy     *uint64
	ProcessedAt     *time.Time
}

func SaveWalletTransaction(walletTransaction *SaveWalletTransactionRequest) error {
	newWalletTransaction := &entity.WalletTransactions{
		UserId:          walletTransaction.UserId,
		OrderId:         walletTransaction.OrderId,
		Amount:          walletTransaction.Amount,
		FlowType:        walletTransaction.FlowType,
		TransactionType: walletTransaction.TransactionType,
		PaymentMethod:   walletTransaction.PaymentMethod,
		Status:          walletTransaction.Status,
		ReferenceNumber: walletTransaction.ReferenceNumber,
		Note:            walletTransaction.Note,
		AdminNote:       walletTransaction.AdminNote,
		ProcessedBy:     walletTransaction.ProcessedBy,
		ProcessedAt:     walletTransaction.ProcessedAt,
	}

	if err := walletTransaction.DB.Create(newWalletTransaction).Error; err != nil {
		return err
	}

	return nil
}
