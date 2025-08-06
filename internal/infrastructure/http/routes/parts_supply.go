package routes

import (
	repository "mecanica_xpto/internal/domain/repository/parts_supply"
	use_case "mecanica_xpto/internal/domain/usecase"

	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func addPartsSupplyRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	partsSupplyUseCase := use_case.NewPartsSupplyUseCase(repository.NewPartsSupplyRepository(db))
	partsSupplyHandler := http.NewPartsSupplyHandler(partsSupplyUseCase)

	partsSupply := rg.Group(PathPartsSupply)
	{
		partsSupply.GET("/:id", partsSupplyHandler.GetPartsSupplyByID)
		partsSupply.GET("/", partsSupplyHandler.ListPartsSupplies)
		partsSupply.POST("/", partsSupplyHandler.CreatePartsSupply)
		partsSupply.PUT("/:id", partsSupplyHandler.UpdatePartsSupply)
		partsSupply.DELETE("/:id", partsSupplyHandler.DeletePartsSupply)
	}
}
