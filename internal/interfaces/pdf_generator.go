package interfaces

import "seblak-bombom-restful-api/internal/model"

type PDFGenerator interface {
	GeneratePDFFromHTML(mail model.CreatePDF) (*model.Attachment, error)
}
