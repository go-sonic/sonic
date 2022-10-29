package dto

type Statistic struct {
	PostCount     int64 `json:"postCount"`
	CommentCount  int64 `json:"commentCount"`
	CategoryCount int64 `json:"categoryCount"`
	TagCount      int64 `json:"tagCount"`
	JournalCount  int64 `json:"journalCount"`
	Birthday      int64 `json:"birthday"`
	EstablishDays int64 `json:"establishDays"`
	LinkCount     int64 `json:"linkCount"`
	VisitCount    int64 `json:"visitCount"`
	LikeCount     int64 `json:"likeCount"`
}
type StatisticWithUser struct {
	Statistic
	User User `json:"user"`
}
