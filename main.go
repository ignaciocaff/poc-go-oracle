package main

import (
	"context"
	"fmt"
	"poc/internal/core"
	"poc/internal/core/env"
	"time"
)

func main() {
	ctx := context.Background()
	config := env.GetEnv(".env.development")

	/*oracleWithSql := core.Oracle{}
	oracleWithSql.OpenOracle(config)
	go oracleWithSql.ExecuteSPWithCursor(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO", nil, "20352579972")

	oracleWithSqlx := core.OracleSqlx{}
	oracleWithSqlx.OpenOracle(config)
	oracleWithSqlx.ExecuteSPWithCursor(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO", nil, "20352579972")
	*/
	/*oracleWithSqlxStament := core.OracleSqlxStatement{}
	oracleWithSqlxStament.OpenOracle(config)
	oracleWithSqlxStament.ExecuteSPWithCursor(ctx)*/

	/*oracleSqlxStatementGorm := core.OracleSqlxStatementGorm{}
	oracleSqlxStatementGorm.OpenOracle(config)
	oracleSqlxStatementGorm.ExecuteSPWithCursor()*/
	/*ctx := context.Background()
	oracleSqlxStatementWithGodror := core.OracleSqlxStatementGodror{}
	db := oracleSqlxStatementWithGodror.OpenOracle(config)
	oracleSqlxStatementWithGodror.ExecuteSPWithCursor(ctx, db)*/

	var res Res

	workingExecution := core.WorkingExecution{}
	workingExecution.OpenOracle(config)
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO", &res, "20352579972")
	fmt.Printf("Desde el main %+v\n", res)
}

type Res struct {
	Cuil            string    `oracle:"Cuit/cuil"`
	Apellido        string    `oracle:"Apellido"`
	Nombre          string    `oracle:"Nombre"`
	FechaDefuncion  time.Time `oracle:"FechaFallecimiento"`
	FechaNacimiento time.Time `oracle:"FechaNacimiento"`
	Genero          bool      `oracle:"Genero, convert"`
	IdLocalidad     int       `oracle:"IdLocalidad"`
	Mail            string    `oracle:"Mail"`
}
