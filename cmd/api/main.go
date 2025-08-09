package main

import (
	"mecanica_xpto/internal/infrastructure/http/routes"
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
	routes.Run()
}
