// gorm gen configure
package main

import (
	"byteurl/dal"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(dal.ConnectMySQL())

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
