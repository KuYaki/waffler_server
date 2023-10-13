package storage

import "gorm.io/gorm"

func ApplyOrder(conn *gorm.DB, orders []string) (tx *gorm.DB) {
	tx = conn
	for _, order := range orders {
		tx = tx.Order(order)
	}
	return
}
