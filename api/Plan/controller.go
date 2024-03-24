package plan

import (
	"gmc-blog-server/model"
	"net/http"
	"strconv"

	plan "gmc-blog-server/view/Plan"
	"log"

	"github.com/gin-gonic/gin"
)

func Submit(c *gin.Context) error {
	var model model.Plan

	err := c.ShouldBind(&model)
	if err != nil {
		return err
	}
	if model.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "请携带用户id",
		})
	}

	log.Printf("plan to add is: %s", model.Content)

	p, err := plan.Search(int(model.UserId))
	log.Printf("plan found id is: %v", p.ID)

	if err != nil {
		return nil
	}

	if p.ID == 0 {
		err = plan.InsertPlan(model)
		if err != nil {
			return err
		}
	} else {
		err = plan.Update(model)
		if err != nil {
			return err
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "计划内容更新成功",
		"data":    model,
	})

	return nil
}

func Search(c *gin.Context) error {
	userIdstr := c.Param("userId")
	userId, err := strconv.Atoi(userIdstr)
	if err != nil {
		return err
	}

	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "请携带用户id",
		})
	}

	p, err := plan.Search(int(userId))

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "查询计划成功",
		"data":    p,
	})

	return nil
}
