package core

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Sqlx struct {
	db *sqlx.DB
}

func (o *Sqlx) ExecuteSPWithCursor(spName string, spResult interface{}, args ...interface{}) error {
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
	var cursor sql.Rows
	execArgs := make([]interface{}, len(args)+1)
	execArgs[0] = sql.Out{Dest: &cursor}
	copy(execArgs[1:], args)

	_, err := o.db.Exec(cmdText, execArgs...)
	if err != nil {
		panic(fmt.Errorf("error executing stored procedure: %w", err))
	}

	if err != nil {
		fmt.Println(err) // column count mismatch: we have 10 columns, but given 0 destination
	}
	defer cursor.Close()
	return nil
}
