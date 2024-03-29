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
	dr := db.DB.GetDbW()
	err := dr.Model(&model.Articles{}).Where("id = ?", article.ID).Updates(map[string]interface{}{
		"title":   article.Title,
		"imgUrl":  article.ImgUrl,
		"type":    article.Type,
		"brief":   article.Brief,
		"content": article.Content,
	}).Error
	return err
}

func ArticleSimplaeInfosQueryByUserId(userId string) ([]model.Articles, error) {
	dr := db.DB.GetDbR()

	log.Println("user id: ", userId)

	var articleSI []model.Articles
	err := dr.Select("id, userId, imgUrl, title, brief, created_at").Where("userId = ?", userId).Order("created_at desc").Find(&articleSI).Error
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
	err := dr.Select("id, userId, created_at, title, imgUrl, brief, content, type").Where("id = ?", articleId).First(&a).Error
	if err != nil {
		return ArticleInfoModel{}, err
	}

	return ArticleInfoModel{
		ID:        int64(a.ID),
		UserId:    a.UserId,
		CreatedAt: a.CreatedAt,
		Content:   a.Content,
		Type:      a.Type,
		Title:     a.Title,
		Brief:     a.Brief,
		ImgUrl:    config.PHOTO_QUERY_PATH + strconv.Itoa(int(a.UserId)) + "/article/" + a.ImgUrl,
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

func ArticleQueryByUserIdSimplaeLife(userId string) (model.Articles, error) {
	dr := db.DB.GetDbR()

	log.Println("user id: ", userId)

	var articleSI model.Articles
	err := dr.Select("id, userId, imgUrl, title, brief, created_at").Where("userId = ?", userId).Order("created_at desc").First(&articleSI).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("未找到数据")
		} else {
			return model.Articles{}, err
		}
	}

	return articleSI, nil
}

func GetArticleInfos(a []model.Articles, userid string) []ArticleSimpleInfoModel {
	var articleSIs []ArticleSimpleInfoModel //id, userId, imgUrl, title, content, created_at
	if len(a) > 0 {
		for _, a := range a {
			articleSIs = append(articleSIs, ArticleSimpleInfoModel{
				ID:     int64(a.ID),
				Title:  a.Title,
				UserId: a.UserId,
				ImgUrl: config.PHOTO_QUERY_PATH + userid + "/article/" + a.ImgUrl,
				Brief:  a.Brief,
				Date:   a.CreatedAt,
			})
		}
	}
	return articleSIs
}
