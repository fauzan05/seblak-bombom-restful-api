package model

import (
	"strings"
)

type Attachment struct {
	Filename string // contoh: "laporan.pdf"
	MimeType string // contoh: "application/pdf"
	Content  []byte // isi file dalam bentuk []byte
}

type Mail struct {
	To          []string
	Cc          []string
	Subject     string
	Template    strings.Builder
	Attachments []Attachment
}
