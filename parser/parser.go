package parser

import (
	"regexp"
	"strings"

	"github.com/tealeg/xlsx"
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

// ParseRow получает строку xlsx, парсит и отправляет заявку через канал
func (p *Parser) ParseRow(r *xlsx.Row, d func(e *Order) bool) int64 {
	extID, err := r.Cells[1].Int64()
	if err != nil {
		return extID
	}
	mku, _ := r.Cells[2].Int()
	mku = mku / 1000
	address := r.Cells[3].Value
	phones := p.Reg.FindAllString(r.Cells[4].Value, -1)
	name := r.Cells[5].Value
	comment := r.Cells[5].Value
	devids := p.ParseDevices(&comment)
	cl := Contact{Name: name, Address: address, Phones: phones}
	o := &Order{ExtID: extID, Mku: int8(mku), Coment: comment, Devices: devids, Client: cl}
	d(o)
	return extID
}

// ParseDevices поиск и замена известный устройств
func (p *Parser) ParseDevices(s *string) []int8 {
	ids := make([]int8, 0)
	for _, d := range p.Devices {
		if strings.Contains(*s, d.Model) {
			ids = append(ids, d.ID)
			*s = strings.Replace(*s, d.Model, " ", -1)
		}
	}
	return ids
}
