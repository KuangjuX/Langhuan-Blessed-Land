package Models

import(
	"time"

	orm "github.com/KuangjuX/Lang-Huan-Blessed-Land/Databases"
)

type Article struct{
	ID 		int64		`json:"id" gorm:"primaryKey" form:"id"`
	UserID  int64		`json:"user_id" form:"user_id"`
	TagID	int64		`json:"tag_id" form:"tag_id"`
	Title	string		`json:"title"`
	Content	string		`json:"content"`
	Created time.Time	`json:"created_at"`
	Update	time.Time	`json:"update_at"`
}




func GetAllArticles() ([]Article, error){
	var articles []Article
	result := orm.Db.Find(&articles)
	err := result.Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func GetArticlesByTag(tag_id string) ([]Article, error){
	var articles []Article
	result := orm.Db.Where("tag_id = ?", tag_id).Find(&articles)
	err := result.Error

	if err != nil{
		return articles, err
	}

	return articles, err
}

func GetArticlesByPage(page, size int, params map[string]string) (interface{}, error){
	DB := orm.Db
	articles := make([]Article, 0)
	if user_id, ok := params["user_id"]; ok == true {
		DB = DB.Where("user_id = ?", user_id)
	}

	if tag_id, ok := params["tag_id"]; ok == true {
		DB = DB.Where("tag_id = ?", tag_id)
	}

	if page > 0 && size > 0 {
		DB = DB.Limit(size).Offset((page - 1)*size)
	}

	if err := DB.Find(&articles).Error; err != nil{
		return nil, err
	}

	return articles, nil
}

func UploadArticle(user_id, tag_id int64, title, content string) (error) {
	article := Article{
		UserID: user_id,
		TagID: tag_id,
		Title: title,
		Content: content,
	}

	DB := orm.Db
	DB = DB.Create(&article)
	if err := DB.Error; err != nil{
		return err
	}
	return nil
}

func ModifyArticle(article_id, tag_id int64, title, content string) (error){
	DB := orm.Db
	DB = DB.Model(Article{}).Where("id = ?", article_id).Updates(Article{
		Title: title,
		TagID: tag_id,
		Content: content })
	if DB.Error != nil {
		return DB.Error
	}
	return nil
}

func DeleteArticle(article_id int64) (error) {
	DB := orm.Db
	DB = DB.Delete(&Article{}, article_id)

	if DB.Error != nil {
		return DB.Error
	}
	return nil
}

func SearchArticles(search_text string) (interface{}, error) {
	DB := orm.Db
	articles := make([]Article, 0)
	result := DB.Where("title LIKE ? or content LIKE ?","%" + search_text + "%", "%" + search_text + "%").Find(&articles)

	if err := result.Error; err != nil{
		return nil, err
	}

	return articles, nil
}