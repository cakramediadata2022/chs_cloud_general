package db_query

import (
	"fmt"

	"gorm.io/gorm"
)

func TruncateTable(DB *gorm.DB, TableName string) error {
	if err := DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", TableName)).Error; err != nil {
		return err
	}
	return nil
}
