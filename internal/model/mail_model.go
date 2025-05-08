package model

import (
	"strings"
)

type Mail struct {
	To       []string
	Cc       []string
	Subject  string
	Template strings.Builder
}
