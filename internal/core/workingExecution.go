package core

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"poc/internal/core/env"
	"reflect"
	"strings"
	"time"

	_ "github.com/godror/godror" // Driver de Oracle para sqlx
	"github.com/jmoiron/sqlx"
)

type WorkingExecution struct {
	db *sqlx.DB
}

type Persona struct {
	Cuil            string    `oracle:"Cuit/cuil"`
	Apellido        string    `oracle:"Apellido"`
	Nombre          string    `oracle:"Nombre"`
	FechaDefuncion  time.Time `oracle:"FechaFallecimiento"`
	FechaNacimiento time.Time `oracle:"FechaNacimiento"`
	Genero          string    `oracle:"Genero"`
	IdLocalidad     int       `oracle:"IdLocalidad"`
	Mail            string    `oracle:"Mail"`
}

func (o *WorkingExecution) OpenOracle(config env.EnvApp) {
	timeZone := "UTC"
	connectionString := fmt.Sprintf(`user="%s" password="%s" timezone="%s" connectString="%s"`, config.DB_USERNAME, config.DB_PASSWORD, timeZone, fmt.Sprintf("%s:%s/%s", config.DB_HOST, config.DB_PORT, config.DB_SERVICE))
	db, err := sqlx.Open("godror", connectionString)
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
	fmt.Println("the time in the database ", queryResultColumnOne)
}

func (o *WorkingExecution) ExecuteStoreProcedure(ctx context.Context, spName string, results interface{}, args ...interface{}) (interface{}, error) {
	conn, err := o.db.Conn(ctx)
	if err != nil {
		log.Printf("error getting connection: %+v", err)
		return nil, err
	}

	var cursor driver.Rows

	cmdText := buildCmdText(spName, args...)

	execArgs := buildExecutionArguments(&cursor, args...)

	if _, err := conn.ExecContext(ctx, cmdText, execArgs...); err != nil {
		log.Printf("error running %q: %+v", cmdText, err)
	}

	cols := cursor.(driver.RowsColumnTypeScanType).Columns()
	rows := make([]driver.Value, len(cols))

	allRows, err := populateRows(cursor, cols, rows)
	if err != nil {
		return nil, err
	}

	mapToSlice(results, cols, allRows)
	cursor.Close()
	return nil, nil
}

func populateRows(cursor driver.Rows, cols []string, rows []driver.Value) ([][]driver.Value, error) {
	var allRows [][]driver.Value
	for {
		if err := cursor.Next(rows); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		newRow := make([]driver.Value, len(rows))
		copy(newRow, rows)
		allRows = append(allRows, newRow)
	}
	return allRows, nil
}

func mapToSlice(slicePtr interface{}, cols []string, allRows [][]driver.Value) error {
	slicePtrValue := reflect.ValueOf(slicePtr)

	if slicePtrValue.Kind() != reflect.Ptr || slicePtrValue.Elem().Kind() != reflect.Slice {
		fmt.Println("it is not a slice")
	}

	sliceType := slicePtrValue.Elem().Type()
	elemType := sliceType.Elem()

	for _, val := range allRows {
		if val != nil {
			newElem := reflect.New(elemType).Elem()
			mapTo(newElem.Addr().Interface(), cols, val)
			slicePtrValue.Elem().Set(reflect.Append(slicePtrValue.Elem(), newElem))

		}
	}
	return nil
}

func mapTo(obj interface{}, cols []string, dests []driver.Value) {
	type CustomMap struct {
		string
		bool
	}
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	tags := make(map[string]CustomMap)

	if v.Kind() != reflect.Struct {
		fmt.Println("it is not a struct")
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		arrayTags := field.Tag.Get("oracle")
		parts := strings.Split(arrayTags, ",")
		tagValue := parts[0]
		convertible := len(parts) > 1 && parts[1] == "convert"
		if tagValue != "" {
			tags[tagValue] = CustomMap{fieldName, convertible}
		}
	}
	for i, col := range cols {
		fieldName := tags[col].string
		field := v.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			fieldType := field.Type()
			val := dests[i]
			if val != nil {
				if tags[col].bool && fieldType.Kind() == reflect.Bool {
					val = val == "S"
				}
				destType := reflect.TypeOf(val)
				if destType.ConvertibleTo(fieldType) {
					field.Set(reflect.ValueOf(val).Convert(fieldType))
				} else {
					fmt.Printf("can not convert %v to %v\n", destType, fieldType)
				}
			} else {
				field.Set(reflect.Zero(fieldType))
			}
		}
	}
}

func buildExecutionArguments(cursor *driver.Rows, args ...interface{}) []interface{} {
	execArgs := make([]interface{}, len(args)+1)
	execArgs[0] = sql.Out{Dest: cursor}
	copy(execArgs[1:], args)
	return execArgs
}

func buildCmdText(spName string, args ...interface{}) string {
	cmdText := fmt.Sprintf("BEGIN %s(:1", spName)
	for i := 0; i < len(args); i++ {
		cmdText += fmt.Sprintf(", :%d", i+2)
	}
	cmdText += "); END;"
	return cmdText
}


func populateRowsOneResult(cursor driver.Rows, cols []string, rows []driver.Value) error {
	for {
		if err := cursor.Next(rows); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}