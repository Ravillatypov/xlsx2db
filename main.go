package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"regexp"

	DB "github.com/Ravillatypov/xlsx2db/db"
	"github.com/Ravillatypov/xlsx2db/parser"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
)

func main() {
	var filename string
	var mysqlconf string
	flag.StringVar(&filename, "f", "/tmp/file.xlsx", "путь к файлу xlsx/xls")
	flag.StringVar(&mysqlconf, "m", "/mydb", "путь к mysql")
	flag.Parse()
	mdb, err := sql.Open("mysql", mysqlconf)
	defer mdb.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-2)
	}
	T := DB.Db{}
	T.Init(mdb)
	devs, err := T.GetDevices()
	if err != nil {
		os.Exit(-3)
	}
	re, _ := regexp.Compile(`[0-9]{11}`)
	P := parser.Parser{Devices: devs, Reg: *re}
	xlsFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	sheet := xlsFile.Sheets[0]
	for _, r := range sheet.Rows {
		P.ParseRow(r, T.Insert)
	}
}
