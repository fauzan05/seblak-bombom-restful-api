package log

import (
	"errors"
	"fmt"
	"seblak-bombom-restful-api/internal/helper"
	"testing"
)


func Pembagian(x, y int) (int, error) {
    if y == 0 {
        return 0, errors.New("tidak dapat membagi dengan nol")
    }
    result := x / y
    return result, nil
}

func TestSaveToLog(t *testing.T) {
	hasil,err := Pembagian(2,0)
	helper.HandleErrorWithPanic(err)
	fmt.Println(hasil)
	// errornya := "ada error"
	// helper.SaveToLogError(errornya)
}