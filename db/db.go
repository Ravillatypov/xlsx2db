package db

import (
	"database/sql"
	"fmt"

	"github.com/Ravillatypov/xlsx2db/parser"
)

// Db - type for write to db
type Db struct {
	db       *sql.DB
	scontact string
	icontact string
	idevice  string
	sdevices string
	sorder   string
	iorder   string
	iphone   string
}

// Init - initialise Db
func (t *Db) Init(d *sql.DB) {
	t.db = d
	t.icontact = `INSERT INTO suz_contacts (id,address,fio) VALUES (NULL, ?, ?)`
	t.sorder = `SELECT id FROM suz_device_orders WHERE ext_id=?`
	t.scontact = `SELECT id FROM suz_contacts WHERE address=? AND fio=?`
	t.iorder = `INSERT INTO suz_device_orders (id,status_id,executor_id,ext_id,contact_id,mku,comment) VALUES (NULL,1,0,?,?,?,?)`
	t.iphone = `INSERT INTO suz_contact_phones (id,contact_id,phone) VALUES ?`
	t.idevice = `INSERT INTO suz_devices_per_order (id,order_id,device_id) VALUES ?`
	t.sdevices = `SELECT id,model FROM suz_devices ORDER BY LENGTH(model) DESC`
}

// Insert - write order to db
func (t *Db) Insert(o *parser.Order) bool {
	if err := t.db.Ping(); err != nil {
		fmt.Println(err.Error())
		return false
	}
	var ext, cid, oid int64
	var sqlstr string
	t.db.QueryRow(t.sorder, o.ExtID).Scan(&ext)
	fmt.Println("ORD_ID: ", ext)
	if ext > 0 {
		return false
	}
	res, err := t.db.Exec(t.icontact, o.Client.Address, o.Client.Name)
	if err != nil {
		t.db.QueryRow(t.scontact, o.Client.Address, o.Client.Name).Scan(&cid)
	} else {
		cid, _ = res.LastInsertId()
	}
	if cid == 0 {
		return false
	}
	fmt.Println("clientID: ", cid)
	res, err = t.db.Exec(t.iorder, o.ExtID, cid, o.Mku, o.Coment)
	if err != nil {
		t.db.QueryRow(t.sorder, o.ExtID).Scan(&oid)
	} else {
		oid, _ = res.LastInsertId()
	}
	if oid == 0 {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println("OrderID: ", oid)
	for i, phone := range o.Client.Phones {
		if i == 0 {
			sqlstr = fmt.Sprintf("(NULL,%d,'%s')", cid, phone)
		} else {
			sqlstr = fmt.Sprintf("%s,(NULL,%d,'%s')", sqlstr, cid, phone)
		}
	}
	fmt.Println("SQL_VALUES: ", sqlstr)
	t.db.Exec(t.iphone, sqlstr)
	for i, dev := range o.Devices {
		if i == 0 {
			sqlstr = fmt.Sprintf("(NULL,%d,%d)", oid, dev)
		} else {
			sqlstr = fmt.Sprintf("%s,(NULL,%d,%d)", sqlstr, oid, dev)
		}
	}
	fmt.Println("SQL_VALUES: ", sqlstr)
	t.db.Exec(t.idevice, sqlstr)
	return true
}

// Run - run writer
func (t *Db) Run(ch chan *parser.Order) {
	select {
	case o := <-ch:
		fmt.Printf("row: ext: %d\tname: %s,\taddress:\t%s,\tcomment: %s\n", o.ExtID, o.Client.Name, o.Client.Address, o.Coment)
		_ = t.Insert(o)
	default:
		fmt.Println("wait...")
	}
}

// GetDevices - get devices list
func (t *Db) GetDevices() ([]parser.Device, error) {
	var devs []parser.Device
	var id int8
	var mod string
	if err := t.db.Ping(); err != nil {
		fmt.Println(err.Error())
		return devs, err
	}

	sdevices, _ := t.db.Prepare(t.sdevices)
	defer sdevices.Close()
	err := t.db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		return devs, err
	}
	res, err := sdevices.Query()
	defer res.Close()
	if err != nil {
		fmt.Println(err.Error())
		return devs, err
	}
	for res.Next() {
		err = res.Scan(&id, &mod)
		if err != nil {
			fmt.Println(err.Error())
		}
		devs = append(devs, parser.Device{ID: id, Model: mod})
	}
	return devs, nil

}
