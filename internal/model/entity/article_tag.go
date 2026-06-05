package entity

// ArticleTag 文章标签关联实体
type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "article_tag"
}
