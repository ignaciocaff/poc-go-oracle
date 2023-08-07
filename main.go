package main

import (
	"context"
	"fmt"

	//"fmt"
	//"github.com/ignaciocaff/oraclesp"
	"poc/internal/core"
	"poc/internal/core/env"
	"time"
)

func main() {
	ctx := context.Background()
	config := env.GetEnv(".env.development")

	var a2 Persona
	var a4 Mensaje
	var a5 Res5
	var a6 []Res4
	var a3 Usuario
	var a1 Res2

	workingExecution := core.WorkingExecution{}
	workingExecution.OpenOracle(ctx, config)
	/*workingExecution.ExecuteStoreProcedure(ctx, "PKG_SEGURIDAD_PERSONAS.PR_OBTENER_DATOS_COMPARAR", &a2, "20352579972")
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_SEGURIDAD_CONSULTAS.PR_OBTENER_USUARIO_X_PERS", &a3, "20352579972", "01", "35257997", 0)
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_SEGURIDAD_PERSONAS.PR_ACTUALIZAR_FUNC_X_CUIL", &a4, "20352579972")
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_SEGURIDAD_CONSULTAS.PR_OBT_TIPOS_USU_X_USR", &a5, 213, 2)
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_TRAMITES_CONSULTAS.PR_OBT_BANCOS", &a6)*/
	workingExecution.ExecuteStoreProcedure(ctx, "PKG_TRAMITES_CONSULTAS.PR_VALIDAR_USUARIO_TRAMITE", &a1, "20213984544", 76160)

	fmt.Printf("%+v\n", a1)
	fmt.Printf("%+v\n", a3)
	fmt.Printf("%+v\n", a2)
	fmt.Printf("%+v\n", a4)
	fmt.Printf("%+v\n", a5)
	fmt.Printf("%+v\n", a6)

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
	TieneTramite bool
}

type Res1 struct {
	NombreTramite string `oracle:"NombreTramite,convert"`
}

type Res4 struct {
	Id     int    `oracle:"Id"`
	Nombre string `oracle:"Nombre"`
}

type Mensaje struct {
	Mensaje string `oracle:"Mensaje"`
}

type Res5 struct {
	Id        int    `oracle:"Id"`
	Nombre    string `oracle:"Nombre"`
	ImageName string `oracle:"NombreImagen"`
}

type Persona struct {
	PerNroInt    int    `oracle:"PerNroInt"`
	Celular      string `oracle:"Celular"`
	Telefono     string `oracle:"Telefono"`
	Correo       string `oracle:"Correo"`
	IdLocalidad  int    `oracle:"IdLocalidad"`
	Barrio       string `oracle:"Barrio"`
	Calle        string `oracle:"Calle"`
	Altura       string `oracle:"Altura"`
	Depto        string `oracle:"Depto"`
	Piso         string `oracle:"Piso"`
	CodigoPostal int    `oracle:"CodigoPostal"`
}

type Usuario struct {
	IdUsuario     int
	Apellido      string
	Nombre        string
	NroDocumento  string
	Cuil          string
	IdTipoUsuario int
	FecAlta       time.Time
}
