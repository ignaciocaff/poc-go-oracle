package main

import (
	//"context"
	"poc/internal/core"
	"poc/internal/core/env")

func main() {
	//ctx := context.Background()
	config := env.GetEnv(".env.development")

	/*oracleWithSql := core.Oracle{}
	oracleWithSql.OpenOracle(config)
	go oracleWithSql.ExecuteSPWithCursor(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO", nil, "20352579972")

	oracleWithSqlx := core.OracleSqlx{}
	oracleWithSqlx.OpenOracle(config)
	oracleWithSqlx.ExecuteSPWithCursor(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO", nil, "20352579972")

	oracleWithSqlxStament := core.OracleSqlxStatement{}
	oracleWithSqlxStament.OpenOracle(config)
	oracleWithSqlxStament.ExecuteSPWithCursor()*/

	/*oracleSqlxStatementGorm := core.OracleSqlxStatementGorm{}
	oracleSqlxStatementGorm.OpenOracle(config)
	oracleSqlxStatementGorm.ExecuteSPWithCursor()*/

	oracleSqlxStatementWithGodror := core.OracleSqlxStatementGodror{}
	oracleSqlxStatementWithGodror.OpenOracle(config)
	oracleSqlxStatementWithGodror.ExecuteSPWithCursor()
}
