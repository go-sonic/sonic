package vo

type ArchiveYear struct {
	Year  int     `json:"year"`
	Posts []*Post `json:"posts"`
}

type ArchiveMonth struct {
	Month int     `json:"month"`
	Posts []*Post `json:"posts"`
}
