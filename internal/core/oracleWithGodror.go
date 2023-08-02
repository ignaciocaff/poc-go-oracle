package core

import (
	"database/sql"
	"fmt"
	"log"
	"poc/internal/core/env"

	"github.com/jmoiron/sqlx"
)

/* Godror*/

type OracleSqlxStatementGodror struct {
	db *sqlx.DB
}

func (o *OracleSqlxStatementGodror) OpenOracle(config env.EnvApp) {
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

func (o *OracleSqlxStatementGodror) ExecuteSPWithCursor() error {
	// Prepare the statement with cursor output
	defer o.db.Close()
	// Definir el nombre del procedimiento almacenado y sus argumentos
	spName := "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO"
	cuil := "20352579972" // Ejemplo de valor de entrada

	// Preparar el comando con el cursor de salida
	cmdText := fmt.Sprintf("BEGIN %s(:1, :2); END;", spName)

	// Definir una variable para almacenar el cursor de salida
	var cursor *sqlx.Rows

	// Ejecutar el procedimiento almacenado con el cursor de salida
	if err := o.db.Select(&cursor, cmdText, cuil, sql.Out{}); err != nil {
		log.Fatalln(err)
	}
	defer cursor.Close()

	// Aquí puedes iterar a través de las filas del cursor y procesar los resultados
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
	}

	return nil
}
