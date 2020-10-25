package main

import (
	"battleship/cmd"
	_ "battleship/docs"
)

// @title Battleship API
// @version 1.0
// @description Battleship API
// @termsOfService http://swagger.io/terms/
// @contact.name Mahmoud AllamehAmiri
// @contact.email m.allamehamiri@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
