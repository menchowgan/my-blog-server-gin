package db

import "gmc-blog-server/model"

func InitUserDB() error {
	dw := DB.GetDbW()

	err := dw.Migrator().CreateTable(&model.User{})

	// if has {
	// 	return nil
	// }

	// err := dw.AutoMigrate(&model.User{})

	if err == nil {
		return nil
	}

	return err
}
