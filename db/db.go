package db

import (
	"database/sql"
	"fmt"

	"github.com/Ravillatypov/xlsx2db/parser"
)

// Db - type for write to db
type Db struct {
	db       *sql.DB
	scontact *sql.Stmt
	icontact *sql.Stmt
	idevice  *sql.Stmt
	sorder   *sql.Stmt
	iorder   *sql.Stmt
	iphone   *sql.Stmt
}

// Init - initialise Db
func (t *Db) Init() {
	sql := `INSERT INTO suz_contacts (id,address,fio) VALUES (NULL, ?, ?)`
	t.icontact, _ = t.db.Prepare(sql)
	sql = `SELECT id FROM suz_device_orders WHERE ext_id=?`
	t.sorder, _ = t.db.Prepare(sql)
	sql = `SELECT id FROM suz_contacts WHERE address='?' AND fio='?' `
	t.scontact, _ = t.db.Prepare(sql)
	sql = `INSERT INTO suz_device_orders (id,ext_id,contact_id,mku,comment) VALUES (NULL,?,?,?,'?')`
	t.iorder, _ = t.db.Prepare(sql)
	sql = `INSERT INTO suz_contact_phones (id,contact_id,phone) VALUES ?`
	t.iphone, _ = t.db.Prepare(sql)
	sql = `INSERT INTO suz_devices_per_order (id,order_id,device_id) VALUES ?`
	t.idevice, _ = t.db.Prepare(sql)
}

// Insert - write order to db
func (t *Db) Insert(o *parser.Order) bool {
	var ext, cid, oid int64
	t.sorder.QueryRow(o.ExtID).Scan(&ext)
	if ext > 0 {
		return false
	}
	_, err := t.icontact.Exec(o.Client.Address, o.Client.Name)
	if err != nil {
		return false
	}
	err = t.scontact.QueryRow(o.Client.Address, o.Client.Name).Scan(&cid)
	if err != nil {
		return false
	}
	_, err = t.iorder.Exec(o.ExtID, cid, o.Mku, o.Coment)
	if err != nil {
		return false
	}
	t.sorder.QueryRow(o.ExtID).Scan(&oid)
	if oid == 0 {
		return false
	}
	sql := ""
	for i, phone := range o.Client.Phones {
		if i == 0 {
			sql = fmt.Sprintf("(NULL,%d,'%s')", cid, phone)
		} else {
			sql += fmt.Sprintf(",(NULL,%d,'%s')", cid, phone)
		}
	}
	t.iphone.Exec(sql)
	for i, dev := range o.Devices {
		if i == 0 {
			sql = fmt.Sprintf("(NULL,%d,%d)", oid, dev)
		} else {
			sql += fmt.Sprintf(",(NULL,%d,%d)", oid, dev)
		}
	}
	t.idevice.Exec(sql)
	return true
}

// Run - run writer
func (t *Db) Run(ch chan *parser.Order) {
	select {
	case o := <-ch:
		t.Insert(o)
	}
}

// GetDevices - get devices list
func (t *Db) GetDevices() ([]parser.Device, error) {
	sql := `SELECT id,model FROM suz_devices ORDER BY LENGTH(model) DESC`
	devs := make([]parser.Device, 0)
	var dev parser.Device
	res, err := t.db.Query(sql)
	if err != nil {
		return devs, err
	}
	for res.Next() {
		err = res.Scan(&dev)
		if err != nil {
			devs = append(devs, dev)
		}
	}
	return devs, nil

}
