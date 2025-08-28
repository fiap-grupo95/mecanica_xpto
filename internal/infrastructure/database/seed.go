package database

import (
	"fmt"
	"mecanica_xpto/internal/domain/model/dto"
)

const (
	defaultPassword = "jKHN1SmGKuGKyiXhbnaOZg==.0/rdilUJyR5raIXVdOCaX8szZCEUzIpIhYTQIMaLwc8="
	defaultUserType = "admin"
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
				Password: defaultPassword,
				UserType: defaultUserType,
			},
			{
				Email:    "joao@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
			},
			{
				Email:    "joana@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
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

	// Seed status_Service_Order_dtos
	var countSOStatus int64
	db.Model(&dto.ServiceOrderStatusDTO{}).Count(&countSOStatus)
	if countSOStatus == 0 {
		serviceOrderStatus := []dto.ServiceOrderStatusDTO{
			{
				Description: "RECEBIDA",
			},
			{
				Description: "EM DIAGNÓSTICO",
			},
			{
				Description: "AGUARDANDO APROVAÇÃO",
			},
			{
				Description: "APROVADA",
			},
			{
				Description: "REJEITADA",
			},
			{
				Description: "EM EXECUÇÃO",
			},
			{
				Description: "FINALIZADA",
			},
			{
				Description: "ENTREGUE",
			},
			{
				Description: "CANCELADA",
			},
		}

		if err := db.Create(&serviceOrderStatus).Error; err != nil {
			fmt.Println("Erro ao criar OS Status:", err)
			return
		}
		fmt.Println("Seeded Service Order Status successfully")
	} else {
		fmt.Println("Service Order Status already seeded")
	}
	// Seed status_Aditional_Repair_dtos
	var countARStatus int64
	db.Model(&dto.AdditionalRepairStatusDTO{}).Count(&countARStatus)
	if countARStatus == 0 {
		additionalRepairStatus := []dto.AdditionalRepairStatusDTO{
			{
				Description: "ABERTA",
			},
			{
				Description: "AGUARDANDO APROVAÇÃO",
			},
			{
				Description: "APROVADA",
			},
			{
				Description: "REJEITADA",
			},
		}

		if err := db.Create(&additionalRepairStatus).Error; err != nil {
			fmt.Println("Erro ao criar Additional Repair Status:", err)
			return
		}
		fmt.Println("Seeded Additional Repair Status successfully")
	} else {
		fmt.Println("Additional Repair Status already seeded")
	}
}
