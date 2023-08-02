# poc-go-oracle
"github.com/godror/godror" GODROR -> Eliminado (Usa el driver ODPI-C que usamos en node por ejemplo, pero para los proyectos Go tenes que compilar de una manera especial y además necesitas si o si el instantclient y que el proyecto tenga acceso al mismo)


Ya sea que use: "github.com/jmoiron/sqlx" o "database/sql" por si solos ambos explotan de la misma forma Falta oracle driver?. Por lo tanto si uso "github.com/sijms/go-ora/v2" funcionan y conectan correctamente a la DB, "github.com/sijms/go-ora/v2" existe otro como este que no sea GODROR y funcione?
Con esto establecimos que al ser sqlx una ampliacion de funcionalidades de sql, podríamos usar sqlx sin problema, Excepto porque sqlx no tiene sql.Out{Dest: &cursor}  que si tiene sql (No se porque, todavía no lo busque) 


La explicación a porque SQLX no lo tiene. Y también la explicación de porque aun no pude hacer andar sql.Out, como cursor de salida y me estoy viendo obligado a usar ora.RefCursor (Ora pertenece a "github.com/sijms/go-ora/v2")
Aparentemente sería la siguiente:

La ejecucion de procedimientos almacenados con cursores es una operacion especifica del "driver"

Ambas implementaciones ya sea que use:
```go    var cursor *sqlx.Rows
    execArgs := make([]interface{}, len(args)+1)
    execArgs[0] = sql.Out{Dest: &cursor}
    copy(execArgs[1:], args)

    _, err := o.db.Exec(cmdText, execArgs...)```
O
```go      var cursor *sql.Rows
    execArgs := make([]interface{}, len(args)+1)
    execArgs[0] = sql.Out{Dest: &cursor}
    copy(execArgs[1:], args)
    _, err := o.db.Exec(cmdText, execArgs...)```
Nos da lo mismo, sqlx o sql, el error es: **panic: error executing stored procedure: unsupported go type**

