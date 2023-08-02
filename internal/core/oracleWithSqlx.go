package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"poc/internal/core/env"

	"github.com/jmoiron/sqlx"
	ora "github.com/sijms/go-ora/v2"
)

type OracleSqlx struct {
	db *sqlx.DB
}

func (o *OracleSqlx) OpenOracle(config env.EnvApp) {
	connectionString := "oracle://" + config.DB_USERNAME + ":" + config.DB_PASSWORD + "@" + config.DB_HOST + ":" + config.DB_PORT + "/" + config.DB_SERVICE
	db, err := sqlx.Open("oracle", connectionString)
	if err != nil {
		panic(err)
	}
	o.db = db
	var queryResultColumnOne string
	row := o.db.QueryRow("SELECT systimestamp FROM dual")
	err = row.Scan(&queryResultColumnOne)
	if err != nil {
		panic(fmt.Errorf("error scanning db: %w", err))
	}
	fmt.Println("The time in the database ", queryResultColumnOne)
}

func (o *OracleSqlx) ExecuteSPWithCursor(ctx context.Context, spName string, spResult interface{}, args ...interface{}) error {
	// Prepare the statement with cursor output
	defer o.db.Close()
	cmdText := fmt.Sprintf("BEGIN %s(:1", spName)

	// Add placeholders for the dynamic parameters
	for i := 0; i < len(args); i++ {
		cmdText += fmt.Sprintf(", :%d", i+2)
	}
	cmdText += "); END;"
	fmt.Println("Sp completo: " + cmdText)
	var cursor ora.RefCursor
	// Prepare the statement with cursor output
	//var cursor ora.RefCursor
	execArgs := make([]interface{}, len(args)+1)
	execArgs[0] = sql.Out{Dest: &cursor}
	copy(execArgs[1:], args)

	_, err := o.db.Exec(cmdText, execArgs...)
	if err != nil {
		panic(fmt.Errorf("error executing stored procedure: %w", err))
	}
	defer cursor.Close()

	rows, err := cursor.Query()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows.Columns())

	/*// Forma de recorrer el cursor si fuese *sqlx.Rows
	for cursor.Next() {
		// Escanear los valores de la fila en un struct u otra estructura
		var resultado struct {
			Columna1 int
			Columna2 string
		}
		if err := cursor.StructScan(&resultado); err != nil {
			log.Fatalln(err)
		}

		// Realizar acciones con los valores escaneados
		fmt.Printf("Columna1: %d, Columna2: %s\n", resultado.Columna1, resultado.Columna2)
	}*/

	return nil
}
