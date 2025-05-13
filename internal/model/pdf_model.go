package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"strings"
)

type CreatePDF struct {
	Filename    string
	HTML        *strings.Builder
	DPI         uint
	Orientation helper.PDFOrientation // gunakan "Portrait" atau "Landscape"
	PageSize    helper.PDFPageSize    // misalnya: "A4", "Letter"
	FooterText  string
}
