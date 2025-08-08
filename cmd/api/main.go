package main

import (
	_ "mecanica_xpto/docs"
	"mecanica_xpto/internal/infrastructure/http/routes"
)

// @title           Mecanica XPTO API
// @version         1.0
// @description     API for vehicle management system
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

func main() {
	routes.Run()
}
