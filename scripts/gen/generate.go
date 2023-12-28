package main

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gorm.io/gen"
)

const MySQLDSN = ""

var DB *gorm.DB

func init() {
	DB = ConnectDB(MySQLDSN).Debug()
}

var dataMap = map[string]func(gorm.ColumnType) (dataType string){}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../../dal",
		ModelPkgPath: "../../model/entity",
		Mode:         gen.WithDefaultQuery,

		WithUnitTest: false,
		//FieldNullable: true,
		//FieldCoverable:    true,
		FieldWithIndexTag: true,
		//FieldSignable:     true,
		FieldWithTypeTag: true,
	})

	g.WithModelNameStrategy(func(tableName string) (modelName string) {
		if tableName == "meta" {
			return "Meta"
		}
		return DB.NamingStrategy.SchemaName(tableName)
	})
	g.UseDB(DB)

	//g.WithDataTypeMap(dataMap)
	//g.WithJSONTagNameStrategy(func(c string) string { return "-" })
	//
	//g.ApplyBasic(dto.AttachmentDTO{})
	g.ApplyBasic(g.GenerateAllTable()...)
	g.GenerateAllTable()

	g.Execute()
}

func ConnectDB(dsn string) (db *gorm.DB) {
	var err error

	if strings.HasSuffix(dsn, "sqlite.db") {
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	} else {
		db, err = gorm.Open(mysql.Open(dsn))
	}

	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}

	return db
}
