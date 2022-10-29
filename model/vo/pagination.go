package vo

type Pagination struct {
	RainbowPages     []RainbowPage
	PrevPageFullPath string
	NextPageFullPath string
	HasPrev          bool
	HasNext          bool
}
