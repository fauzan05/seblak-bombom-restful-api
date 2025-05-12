package tests

import (
	"fmt"
	"html/template"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/helper/generate_file"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePDF(t *testing.T) {
	newCreatePDF := new(model.CreatePDF)
	newCreatePDF.DPI = 300
	newCreatePDF.Filename = "Invoice-1"
	newCreatePDF.FooterText = "LALALA"
	newCreatePDF.Orientation = helper.PORTRAIT
	newCreatePDF.PageSize = helper.A4
	templatePath := "../internal/templates/pdf/orders/invoice.html"
	tmpl, err := template.ParseFiles(templatePath)
	assert.Nil(t, err)
	
	items := []map[string]string{
		{
			"Name":  "Item A",
			"Quantity":   "1",
			"UnitPrice": "Rp10.000",
			"TotalPrice": "Rp10.000",
		},
		{
			"Name":  "Item B",
			"Quantity":   "2",
			"UnitPrice": "Rp20.000",
			"TotalPrice": "Rp40.000",
		},
	}
	bodyBuilder := new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]any{
		"InvoiceNumber": "INV/asdsa/MPL/423423",
		"PurchaseDate": "24 Maret 2025",
		"SellerName": "Three Acc",
		"BuyerName": "Fauzan Nur hidayat",
		"ShippingAddress": "Sample Address",
		"Items": items,
		"Subtotal": "259.000",
		"Discount": "5.180",
		"ShippingCost": "37.000",
		"TotalShopping": "292.329",
		"TotalBilling": "293.329",
		"ServiceFee": "1.000",
		"PaymentMethod": "Mandiri Virtual Account",
	})
	assert.Nil(t, err)

	newCreatePDF.HTML = bodyBuilder
	generatePdf, err := generate_file.GeneratePDFFromHTML(*newCreatePDF)
	assert.Nil(t, err)

	saveToDir := "../tmp/orders/"
	filePath, err := generate_file.SaveAttachmentToFile(generatePdf, saveToDir)
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf("../tmp/orders/%s", generatePdf.Filename), filePath)
}
