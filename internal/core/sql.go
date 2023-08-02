package core

import "fmt"

type Sql struct { }


func (o *Sql) ExecuteSPWithCursor(spName string, spResult interface{}, args ...interface{}) error{
	fmt.Println("Hello, World Sql!")
	return nil
}