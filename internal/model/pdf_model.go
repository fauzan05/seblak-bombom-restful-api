package model

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"strings"
)

type CreatePDF struct {
	Filename    string
	HTML        *strings.Builder
	DPI         uint
	Orientation enum_state.PDFOrientation // gunakan "Portrait" atau "Landscape"
	PageSize    enum_state.PDFPageSize    // misalnya: "A4", "Letter"
	FooterText  string
}
