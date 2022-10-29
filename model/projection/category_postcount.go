package projection

type CategoryPostCountProjection struct {
	CategoryID int32 `gorm:"column:category_id"`
	PostCount  int32 `gorm:"column:post_count"`
}
