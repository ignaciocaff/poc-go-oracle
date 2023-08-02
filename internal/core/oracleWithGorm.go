package core

import (
	"fmt"
	"log"
	"poc/internal/core/env"
	"strconv"
	"time"

	oracle "github.com/godoes/gorm-oracle"
	cur "github.com/sijms/go-ora/v2"
	"gorm.io/gorm"
)

/* Descartado Cengsin por godror*/

type OracleSqlxStatementGorm struct {
	db *gorm.DB
}

func (o *OracleSqlxStatementGorm) OpenOracle(config env.EnvApp) {
	// connectionString := "oracle://" + config.DB_USERNAME + ":" + config.DB_PASSWORD + "@" + config.DB_HOST + ":" + config.DB_PORT + "/" + config.DB_SERVICE
	// oracle://user:password@127.0.0.1:1521/service
	port, err := strconv.Atoi(config.DB_PORT)
	url := oracle.BuildUrl(config.DB_HOST, port, config.DB_SERVICE, config.DB_USERNAME, config.DB_PASSWORD, nil)
	db, err := gorm.Open(oracle.Open(url), &gorm.Config{})
	//db, err := sqlx.Open("oracle", connectionString)
	if err != nil {
		panic(err)
	}
	o.db = db
	var result time.Time
	row := o.db.Raw("SELECT systimestamp FROM dual").Scan(&result)
	if row.Error != nil {
		log.Fatal(row.Error)
	}
	fmt.Println("The time in the database ", result)
}

func (o *OracleSqlxStatementGorm) ExecuteSPWithCursor() error {
	// Prepare the statement with cursor output
	// Definir el nombre del procedimiento almacenado y sus argumentos
	spName := "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO"
	cuil := "20352579972" // Ejemplo de valor de entrada

	// Preparar el comando con el cursor de salida
	cmdText := fmt.Sprintf("BEGIN %s(:1, :2); END;", spName)
	var cursor cur.RefCursor

	o.db.Raw(cmdText, &cursor, cuil).Scan(&cursor)
	rows, err := cursor.Query()
	if err != nil {
		log.Fatal(err)
	}
	var (
		var1 string
	)
	for rows.Next_() {
		err = rows.Scan(&var1)
		// check for error
		fmt.Println(var1)
	}

	return nil
}
