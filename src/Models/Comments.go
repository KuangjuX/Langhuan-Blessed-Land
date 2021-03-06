package Models

import (
	"time"

	orm "github.com/KuangjuX/Lang-Huan-Blessed-Land/DataBases/mysql"
)

type Comment struct{
	ID			int			`json:"id"`
	UserID		int			`json:"user_id"`
	ArticleID	int			`json:"article_id"`
	ToCommentID	int			`json:"to_comment_id"`
	Content		string		`json:"content"`
	CreatedAt	time.Time	`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt	time.Time 	`json:"updated_at" gorm:"autoUpdateTime"`
}

func GetCommentsByArticle(article_id int)([]Comment, error){
	comments := make([]Comment, 0)
	DB := orm.Db
	result := DB.Where("article_id = ?", article_id).Find(&comments)
	if err := result.Error; err != nil {
		return nil, err
	}

	return comments, nil
}