package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"poc/internal/core/env"

	ora "github.com/sijms/go-ora/v2"
)

type Oracle struct {
	db *sql.DB
}

func (o *Oracle) OpenOracle(config env.EnvApp) {
	connectionString := "oracle://" + config.DB_USERNAME + ":" + config.DB_PASSWORD + "@" + config.DB_HOST + ":" + config.DB_PORT + "/" + config.DB_SERVICE
	db, err := sql.Open("oracle", connectionString)
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

func (o *Oracle) ExecuteSPWithCursor(ctx context.Context, spName string, spResult interface{}, args ...interface{}) error {
	defer o.db.Close()
	// Prepare the statement with cursor output
	cmdText := fmt.Sprintf("BEGIN %s(:1", spName)

	// Add placeholders for the dynamic parameters
	for i := 0; i < len(args); i++ {
		cmdText += fmt.Sprintf(", :%d", i+2)
	}
	cmdText += "); END;"
	fmt.Println("Sp completo: " + cmdText)
	// Prepare the statement with cursor output
	//var cursor ora.RefCursor
	var cursor ora.RefCursor
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
	/*// Forma de recorrer el cursor si fuese sqlx.Rows
	for cursor.Next() {
		// Escanear los valores de la fila en variables
		var columna1 int
		var columna2 string
		if err := cursor.Scan(&columna1, &columna2); err != nil {
			panic(err)
		}

		// Realizar acciones con los valores escaneados
		fmt.Printf("Columna1: %d, Columna2: %s\n", columna1, columna2)
	}*/

	return nil
}
