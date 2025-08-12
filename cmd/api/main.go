package main

import (
	_ "mecanica_xpto/docs"
	"mecanica_xpto/internal/infrastructure/database"
	"mecanica_xpto/internal/infrastructure/http/routes"
	"os"
)

// @title           Mecanica XPTO API
// @version         1.0
// @description     API for Mecanica XPTO service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

// @BasePath  /v1
// @securityDefinitions.basic  BasicAuth

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "seed":
			database.Seed()
			return
		case "migrate":
			database.Migrate()
			return
		}
	}
	routes.Run()
}
