## Poc oracle store procedure cursors 

### Objetivo
Probar si se puede ejecutar un store procedure de oracle que devuelve un cursor y leerlo con Go

### Primeras pruebas

**github.com/godror/godror GODROR** -> Eliminado (Usa el driver ODPI-C que usamos en node por ejemplo, pero para los proyectos Go tenes que compilar de una manera especial y además necesitas si o si el instantclient y que el proyecto tenga acceso al mismo)


Ya sea que use: **github.com/jmoiron/sqlx** o **database/sql** por si solos ambos explotan de la misma forma Falta oracle driver?. Por lo tanto si uso **github.com/sijms/go-ora/v2** funcionan y conectan correctamente a la DB, **github.com/sijms/go-ora/v2** existe otro como este que no sea GODROR y funcione?
Con esto establecimos que al ser sqlx una ampliacion de funcionalidades de sql, podríamos usar sqlx sin problema, Excepto porque sqlx no tiene 
``` go
sql.Out{Dest: &cursor}
```
que si tiene sql 

La explicación a porque SQLX no lo tiene. Y también la explicación de porque aun no pude hacer andar sql.Out, como cursor de salida y me estoy viendo obligado a usar ora.RefCursor (Ora pertenece a **github.com/sijms/go-ora/v2**)
Aparentemente sería la siguiente:

La ejecucion de procedimientos almacenados con cursores es una operacion especifica del "driver"

### Usando `sqlx`
``` go  var cursor *sqlx.Rows
    execArgs := make([]interface{}, len(args)+1)
    execArgs[0] = sql.Out{Dest: &cursor}
    copy(execArgs[1:], args)

    _, err := o.db.Exec(cmdText, execArgs...)
```
### Usando `sql`

``` go      var cursor *sql.Rows
    execArgs := make([]interface{}, len(args)+1)
    execArgs[0] = sql.Out{Dest: &cursor}
    copy(execArgs[1:], args)
    _, err := o.db.Exec(cmdText, execArgs...)
```

Nos da lo mismo, sqlx o sql, el error es: **panic: error executing stored procedure: unsupported go type**



Si **modifico** el tipo del cursor para ambos **casos sqlx y sql** de esta forma:

```go  
var cursor ora.RefCursor -> Recordando que viene del driver "github.com/sijms/go-ora/v2"
```


Funciona la llamada a los SPS y devuelve valores. El problema radica en que en ambos casos también, me veo **obligado a cambiar  la forma de recorrerlo**, que antes era:

### Usando `sql`
```go
 for cursor.Next() {
        // Escanear los valores de la fila en variables
        var columna1 int
        var columna2 string
        if err := cursor.Scan(&columna1, &columna2); err != nil {
            panic(err)
        }

        // Realizar acciones con los valores escaneados
        fmt.Printf("Columna1: %d, Columna2: %s\n", columna1, columna2)
    }
```

### Usando `sqlx`
```go
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
```

### La forma actual que estamos queriendo modificar es:

```go
   rows, err := cursor.Query()
    // check for error

    var (
        var_1 int64
        var_2 string
    )
    for rows.Next_() {
        err = rows.Scan(&var_1, &var_2)
        // check for error
        fmt.Println(var_1, var_2)
    }
```

Si lo quisiera hacer de esta manera  es sqlx pero no haciendo **Exec** sino **Select** el error cambia:
**call register type before use user defined type (UDT) esto esta vinculado al driver github.com/sijms/go-ora/v2 aparentemente**

```go
    spName := "PKG_TRAMITES_CONSULTAS.PR_OBT_PARENTESCOS"
    cuil := "20352579972" 
    cmdText := fmt.Sprintf("BEGIN %s(:1, :2); END;", spName)
    var cursorsqlx.Rows
    if err := o.db.Select(&cursor, cmdText, cuil, sql.Out{}); err != nil {
        log.Fatalln(err)
    }
    defer cursor.Close()
    for cursor.Next() {
        var resultado struct {
            Columna1 int
            Columna2 string
        }
        if err := cursor.StructScan(&resultado); err != nil {
            log.Fatalln(err)
        }
        fmt.Printf("Columna1: %d, Columna2: %s\n", resultado.Columna1, resultado.Columna2)
    }
```

### Drivers descartados:

github.com/cengsin/oracle -> utiliza godror
github.com/godoes/gorm-oracle -> No funciona


### GORM 

```go
o.db.Raw(cmdText, sql.Out{}, cuil).Scan(results)
```
Devuelve el mismo error **call register type before use user defined type (UDT)**


Si hago lo mismo con GORM pero intentando pasar el **RefCursor** de **github.com/sijms/go-ora/v2**"** estoy en la misma, **call register type before use user defined type (UDT)**
```go
o.db.Raw(cmdText, sql.Out{Dest: &cursor}, cuil).Scan(&results)
```

### Conclusiones parciales 

- Nada que no sea usar el cursor de github.com/sijms/go-ora/v2 y por lo tanto esa forma de recorrer funciona hasta acá
- GORM -> call register type before use user defined type (UDT) con cualquiera driver o misma forma de recorrer si paso cursor
- gobuffalo/pop -> usa godror

**Opciones restantes:**
- Buscar la manera de mapear sin que sea generico el Next() y el Scan() -> No se si se puede
- Meterme en como esta compilado godror con gcc en Windows para que me funcione la librería
- Encontrar otro driver que tenga mejor forma de leer en su cursor