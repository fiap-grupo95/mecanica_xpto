package database

import (
	"fmt"
	"mecanica_xpto/internal/domain/model/dto"
)

func Seed() {
	db := ConnectDatabase()

	// Seed user_dtos
	var countUsers int64
	db.Model(&dto.UserDTO{}).Count(&countUsers)
	if countUsers == 0 {
		users := []dto.UserDTO{
			{
				Email:    "admin@xpto.com",
				Password: "admin123",
				UserType: "admin",
			},
			{
				Email:    "joao@xpto.com",
				Password: "admin123",
				UserType: "admin",
			},
			{
				Email:    "joana@xpto.com",
				Password: "admin123",
				UserType: "admin",
			},
		}

		if err := db.Create(&users).Error; err != nil {
			fmt.Println("Erro ao criar usuários:", err)
			return
		}
		fmt.Println("Seeded users successfully")
	} else {
		fmt.Println("Users already seeded")
	}
}
