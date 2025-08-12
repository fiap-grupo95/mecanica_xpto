package database

import (
	"fmt"
	"mecanica_xpto/internal/domain/model/dto"
)

func Migrate() {
	db := ConnectDatabase()

	err := db.AutoMigrate(
		&dto.PartsSupplyDTO{},
		&dto.ServiceDTO{},
		&dto.VehicleDTO{},
		&dto.ServiceOrderDTO{},
		&dto.CustomerDTO{},
		&dto.UserDTO{},
		&dto.ServiceOrderStatusDTO{},
		&dto.AdditionalRepairDTO{},
		&dto.PartsSupplyServiceOrderDTO{},
		&dto.AdditionalRepairStatusDTO{},
		&dto.UserTypeDTO{},
		&dto.ServiceServiceOrderDTO{},
		&dto.PaymentDTO{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database migrated successfully")
}
