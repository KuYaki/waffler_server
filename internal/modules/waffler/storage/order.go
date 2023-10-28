package storage

import "gorm.io/gorm"

func ApplyOrder(conn *gorm.DB, orders []string) (tx *gorm.DB) {
	tx = conn
	for i, order := range orders {
		if i == 0 {
			order += " NULLS LAST"
		}
		tx = tx.Order(order)
	}
	return
}
