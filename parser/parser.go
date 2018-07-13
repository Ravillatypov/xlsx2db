package parser

import (
	"regexp"
)

// Device модель
type Device struct {
	ID    int8   `db:"id"`
	Model string `db:"model"`
}

// Contact Контакт
type Contact struct {
	Name    string
	Address string
	Phones  []string
}

// Order Заявка
type Order struct {
	ExtID   int64
	Mku     int8
	Coment  string
	Devices []int8
	Client  Contact
}

//Parser Парсер
type Parser struct {
	Devices []Device
	Reg     regexp.Regexp
}
