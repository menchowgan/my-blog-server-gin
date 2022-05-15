package article

import (
	"errors"
	"gmc-blog-server/config"
	"gmc-blog-server/db"
	"gmc-blog-server/model"
	"log"
	"strconv"

	"gorm.io/gorm"
)

func InsertArticle(article *model.Articles) error {
	dw := db.DB.GetDbW()

	log.Printf("Id: %v, Title: %s, UserId: %v\n", article.ID, article.Title, article.UserId)
	log.Println(article)

	err := dw.Create(&article).Error
	if err == nil {
		return nil
	}
	return err
}

func Save(article *model.Articles) error {
	return nil
}

func ArticleSimplaeInfosQueryByUserId(userId string) ([]model.Articles, error) {
	dr := db.DB.GetDbR()

	log.Println("user id: ", userId)

	var articleSI []model.Articles
	err := dr.Select("id, userId, imgUrl, title, brief, created_at").Where("userId = ?", userId).Order("created_at desc").Limit(5).Find(&articleSI).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("未找到数据")
		} else {
			return []model.Articles{}, err
		}
	}

	return articleSI, nil
}

func SearchArticleInfo(articleId string) (ArticleInfoModel, error) {
	var a model.Articles

	dr := db.DB.GetDbR()
	err := dr.Select("id, userId, created_at, content").Where("id = ?", articleId).First(&a).Error
	if err != nil {
		return ArticleInfoModel{}, err
	}

	return ArticleInfoModel{
		ID:      int64(a.ID),
		UserId:  a.UserId,
		Date:    a.CreatedAt,
		Content: a.Content,
	}, nil
}

func SearchArticleInfoByType(userid string, atype string) ([]ArticleSimpleInfoModel, error) {
	var articles []model.Articles

	dr := db.DB.GetDbR()
	err := dr.Select("id, userId, imgUrl, title, brief, created_at").Limit(10).Offset(0).Where("userId = ? AND (type LIKE ? OR title LIKE ? OR brief LIKE ?)", userid, "%"+atype+"%", "%"+atype+"%", "%"+atype+"%").Find(&articles).Error
	if err != nil {
		return []ArticleSimpleInfoModel{}, err
	}
	log.Println("aerticles query by type")

	var as []ArticleSimpleInfoModel
	if len(articles) > 0 {
		for _, a := range articles {
			as = append(as, ArticleSimpleInfoModel{
				ID:     int64(a.ID),
				UserId: a.UserId,
				ImgUrl: config.PHOTO_QUERY_PATH + strconv.Itoa(int(a.UserId)) + "/article/" + a.ImgUrl,
				Title:  a.Title,
				Brief:  a.Brief,
				Date:   a.CreatedAt,
			})
		}
	}

	return as, nil
}
