package core

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"

	//"io"
	"poc/internal/core/env"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	go_ora "github.com/sijms/go-ora/v2"
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

func (o *WorkingExecution) OpenOracle(ctx context.Context, config env.EnvApp) *sqlx.DB {
	/*timeZone := "UTC"
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

	return db*/
	urlOptions := map[string]string{}
	port, err := strconv.Atoi(config.DB_PORT)
	databaseUrl := go_ora.BuildUrl(config.DB_HOST, port, config.DB_SERVICE, config.DB_USERNAME, config.DB_PASSWORD, urlOptions)
	fmt.Println("connection string: ", databaseUrl)

	conn, err := sqlx.ConnectContext(ctx, "oracle", databaseUrl) //db, err := sqlx.Open("goracle", connectionString)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}

	var queryResultColumnOne string
	row := conn.QueryRow("SELECT systimestamp FROM dual")
	err = row.Scan(&queryResultColumnOne)
	if err != nil {
		panic(fmt.Errorf("error scanning db: %w", err))
	}
	fmt.Println("The time in the database ", queryResultColumnOne)
	o.db = conn

	return conn
}

func (o *WorkingExecution) ExecuteStoreProcedure(ctx context.Context, spName string, results interface{}, args ...interface{}) error {
	first := time.Now()

	fmt.Println("Starting procedure " + spName + " time " + first.String())

	resultsVal := reflect.ValueOf(results)

	var cursor go_ora.RefCursor
	cmdText := buildCmdText(spName, args...)
	execArgs := buildExecutionArguments(&cursor, args...)

	_, err := o.db.ExecContext(ctx, cmdText, execArgs...)

	if err != nil {
		panic(fmt.Errorf("error scanning db: %w", err))
	}

	rows, err := cursor.Query()
	if err != nil {
		return err
	}
	cols := rows.Columns()
	dests := make([]driver.Value, len(cols))

	if resultsVal.Kind() == reflect.Ptr && resultsVal.Elem().Kind() == reflect.Slice {
		allRows, err := populateRows(rows, cols, dests)
		if err != nil {
			return err
		}
		mapToSlice(results, cols, allRows)
	} else {
		populateOne(rows, cols, dests)
		mapTo(results, cols, dests)
	}
	cursor.Close()
	fmt.Println("Ending procedure " + spName + " time " + time.Now().String())
	return nil
}

func populateRows(cursor *go_ora.DataSet, cols []string, rows []driver.Value) ([][]driver.Value, error) {
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
				if fieldType.Kind() == reflect.String {
					val = trimTrailingWhitespace(val.(string))
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

func buildExecutionArguments(cursor *go_ora.RefCursor, args ...interface{}) []interface{} {
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

func trimTrailingWhitespace(input string) string {
	if len(input) == 0 {
		return input
	}
	input = strings.TrimRight(input, " ")
	return input
}

func populateOne(cursor *go_ora.DataSet, cols []string, rows []driver.Value) error {
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
