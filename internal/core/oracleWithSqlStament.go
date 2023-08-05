package core

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"poc/internal/core/env"
	"time"
	"github.com/jmoiron/sqlx"
)

type OracleSqlxStatement struct {
	db *sqlx.DB
}

type Resultado struct {
	Cuil            string    `oracle:"Cuit/cuil"`
	Apellido        string    `oracle:"Apellido"`
	Nombre          string    `oracle:"Nombre"`
	FechaDefuncion  time.Time `oracle:"FechaFallecimiento"`
	FechaNacimiento time.Time `oracle:"FechaNacimiento"`
	Genero          string    `oracle:"Genero"`
	IdLocalidad     int       `oracle:"IdLocalidad"`
	Mail			string    `oracle:"Mail"`

}

func (o *OracleSqlxStatement) OpenOracle(config env.EnvApp) {
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

func (o *OracleSqlxStatement) ExecuteSPWithCursor(ctx context.Context) error {
	// Prepare the statement with cursor output
	defer o.db.Close()
	// Definir el nombre del procedimiento almacenado y sus argumentos
	const query = `BEGIN PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO(:1, :2); END;`
	cuil := "20352579972" // Ejemplo de valor de entrada
	// Definir una variable para almacenar el cursor de salida
	var cursor driver.Rows
	execArgs := make([]interface{}, 2)
	execArgs[0] = sql.Out{Dest: &cursor}
	execArgs[1] = cuil

	conn, err := o.db.Conn(ctx)
	if err != nil {
		log.Printf("error getting connection: %+v", err)
	}
	// Ejecutar el procedimiento almacenado con el cursor de salida
	if _, err := conn.ExecContext(ctx, query, sql.Out{Dest: &cursor}, cuil); err != nil {
		log.Printf("error running %q: %+v", query, err)
	}
	cols := cursor.(driver.RowsColumnTypeScanType).Columns()
	rows := make([]driver.Value, len(cols))
	for {
		if err := cursor.Next(rows); err != nil {

			if err == io.EOF {
				break
			}
			cursor.Close()
			return err
		}
	}

	var resultado Resultado
	//mapTo(&resultado, cols, rows)
	fmt.Println(cols)
	fmt.Printf("%+v\n", resultado)
	return nil
}

/*func mapTo(obj interface{}, cols []string, dests []driver.Value) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	tags := make(map[string]string)

	if v.Kind() != reflect.Struct {
		fmt.Println("it is not a struct")
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		tagValue := field.Tag.Get("oracle")

		if tagValue != "" {
			tags[tagValue] = fieldName
		}
	}
	for i, col := range cols {
		fieldName := tags[col]
		field := v.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			fieldType := field.Type()
			val := dests[i]
			if val != nil {
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

}*/
