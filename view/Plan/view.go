package plan

import (
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"log"
)

func InsertPlan(plan model.Plan) error {
	dw := db.DB.GetDbW()

	err := dw.Create(&plan).Error
	if err != nil {
		return err
	}

	return nil
}

func Update(plan model.Plan) error {
	dw := db.DB.GetDbW()

	err := dw.Model(&plan).Where("userId = ?", plan.UserId).Updates(
		map[string]interface{}{
			"content": plan.Content,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

func Search(userId int) (PlanModel, error) {
	var m model.Plan

	dr := db.DB.GetDbR()
	err := dr.Where("userId = ?", userId).First(&m).Error

	log.Printf("plan found: %s", m.Content)

	planModel := PlanModel{}
	if err != nil {
		return planModel, nil
	}

	planModel.Content = m.Content
	planModel.ID = m.ID
	planModel.UserId = m.UserId
	return planModel, nil
}
