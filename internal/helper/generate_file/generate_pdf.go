package generate_file

import (
	"fmt"
	"os"
	"path/filepath"
	"seblak-bombom-restful-api/internal/model"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GeneratePDFFromHTML(pdfConfig model.CreatePDF) (*model.Attachment, error) {
	// Inisialisasi PDF Generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF generator : %w", err)
	}

	// Buat halaman dari HTML content
	page := wkhtmltopdf.NewPageReader(strings.NewReader(pdfConfig.HTML.String()))
	page.EnableLocalFileAccess.Set(true)
	
	if pdfConfig.FooterText != "" {
		page.FooterRight.Set(pdfConfig.FooterText)
		page.FooterFontSize.Set(8)
	}

	pdfg.AddPage(page)
	if pdfConfig.DPI > 0 {
		pdfg.Dpi.Set(pdfConfig.DPI)
	} else {
		pdfg.Dpi.Set(300) // default
	}

	if pdfConfig.Orientation != "" {
		pdfg.Orientation.Set(string(pdfConfig.Orientation))
	} else {
		pdfg.Orientation.Set("Portrait")
	}

	if pdfConfig.PageSize != "" {
		pdfg.PageSize.Set(string(pdfConfig.PageSize))
	} else {
		pdfg.PageSize.Set("A4")
	}

	// Generate PDF
	err = pdfg.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Kembalikan sebagai Attachment
	return &model.Attachment{
		Filename: fmt.Sprintf("%s.pdf", pdfConfig.Filename),
		MimeType: "application/pdf",
		Content:  pdfg.Bytes(),
	}, nil
}


func SaveAttachmentToFile(att *model.Attachment, dir string) (string, error) {
	// Pastikan folder tujuan ada
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(dir, att.Filename)

	// Tulis file ke disk
	err = os.WriteFile(filePath, att.Content, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
