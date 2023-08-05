package core

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"poc/internal/core/env"
	"github.com/jmoiron/sqlx"
)

/* Godror*/

type OracleSqlxStatementGodror struct {
	db *sqlx.DB
}

func (o *OracleSqlxStatementGodror) OpenOracle(config env.EnvApp) *sqlx.DB {
	timeZone := "UTC"
	connectionString := fmt.Sprintf(`user="%s" password="%s" timezone="%s" connectString="%s"`, config.DB_USERNAME, config.DB_PASSWORD, timeZone, fmt.Sprintf("%s:%s/%s", config.DB_HOST, config.DB_PORT, config.DB_SERVICE))
	db, err := sqlx.Open("godror", connectionString)
	if err != nil {
		panic(err)
	}
	var queryResultColumnOne string
	row := db.QueryRow("SELECT systimestamp FROM dual")
	err = row.Scan(&queryResultColumnOne)
	if err != nil {
		panic(fmt.Errorf("error scanning db: %w", err))
	}
	fmt.Println("The time in the database ", queryResultColumnOne)

	return db
}

func (o *OracleSqlxStatementGodror) ExecuteSPWithCursor(ctx context.Context, db *sqlx.DB) error {
	// Prepare the statement with cursor output
	defer db.Close()
	// Definir el nombre del procedimiento almacenado y sus argumentos
	cuil := "20352579972" // Ejemplo de valor de entrada
	var rset1 driver.Rows
	const query = `BEGIN PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO(:1, :2); END;`
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Printf("Error getting connection: %+v", err)
	}

	if _, err := conn.ExecContext(ctx, query, sql.Out{Dest: &rset1}, cuil); err != nil {
		log.Printf("Error running %q: %+v", query, err)
	}
	defer rset1.Close()

	cols1 := rset1.(driver.RowsColumnTypeScanType).Columns()
	dests1 := make([]driver.Value, len(cols1))
	for {
		if err := rset1.Next(dests1); err != nil {
			if err == io.EOF {
				break
			}
			rset1.Close()
			return err
		}
		fmt.Println(dests1)
	}

	return nil
}
