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

	var res1 []Res1
	var res2 []Res2
	var res3 []Res3

	workingExecution := core.WorkingExecution{}
	workingExecution.OpenOracle(config)
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_NOMBRE_TIPO_TRAMITE", &res1, 4)
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_TRAMITES_CONSULTAS.PR_VALIDAR_USUARIO_TRAMITE", &res2, "23179917939", 54075)
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO", &res3, "20352579972")

	fmt.Printf("Desde el main %+v\n", res1)
	fmt.Printf("Desde el main %+v\n", res2)
	fmt.Printf("Desde el main %+v\n", res3)
}

type Res3 struct {
	Cuil            string    `oracle:"Cuit/cuil"`
	Apellido        string    `oracle:"Apellido"`
	Nombre          string    `oracle:"Nombre"`
	FechaDefuncion  time.Time `oracle:"FechaFallecimiento"`
	FechaNacimiento time.Time `oracle:"FechaNacimiento"`
	Genero          bool      `oracle:"Genero,convert"`
	IdLocalidad     int       `oracle:"IdLocalidad"`
	Mail            string    `oracle:"Mail"`
}

type Res2 struct {
	HasTramit bool `oracle:"TieneTramite,convert"`
}

type Res1 struct {
	NombreTramite string `oracle:"NombreTramite,convert"`
}
