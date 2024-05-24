package models

import (
	"gorm.io/gorm"
)

type Events struct {
	id       uint    `gorm: "primary key;autoIncrement"	json: "id"`
	ClientId *string `json: "clientId"`
	Type     *string `json: "type"`
	Source   *string `json: "source"`
}

func MigrateEvents(db *gorm.DB) error {
	err := db.AutoMigrate(&Events{})
	return err
}
