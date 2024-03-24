package plan

import (
	"encoding/json"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	r "gmc-blog-server/redis"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
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

	_, err = searchAndUpdate(strconv.Itoa(int(plan.UserId)))
	if err != nil {
		return err
	}

	return nil
}

func Search(userId int) (PlanModel, error) {
	rKey := r.GetKey("plan", strconv.Itoa(userId))
	s, err := r.RedisDb.Get(rKey).Result()
	if err == redis.Nil || err != nil {
		return searchAndUpdate(strconv.Itoa(userId))
	} else {
		planModel := PlanModel{}
		if s != "" {
			log.Println("cached plan", s)
			err := json.Unmarshal([]byte(s), &planModel)
			if err != nil {
				return planModel, err
			}
		}
		return planModel, nil
	}
}

func searchAndUpdate(userId string) (PlanModel, error) {
	var m model.Plan
	planModel := PlanModel{}
	rKey := r.GetKey("plan", userId)
	dr := db.DB.GetDbR()
	err := dr.Where("userId = ?", userId).First(&m).Error
	log.Printf("plan found: %s", m.Content)

	if err != nil {
		return planModel, err
	}

	planModel.Content = m.Content
	planModel.ID = m.ID
	planModel.UserId = m.UserId
	if pl, err := json.Marshal(&planModel); err == nil {
		r.RedisDb.Set(rKey, pl, 96*time.Hour)
	}
	return planModel, nil
}
