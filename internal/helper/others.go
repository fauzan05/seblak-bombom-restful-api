package helper

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
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
