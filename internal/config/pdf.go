package config

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/sirupsen/logrus"
)

func NewPDFGenerator(log *logrus.Logger) *wkhtmltopdf.PDFGenerator {
	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	// if err != nil {
	// 	log.Fatalf("error pdf generator : %v", err)
	// }

	return pdfg
}
