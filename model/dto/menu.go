package dto

type Menu struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Priority int32  `json:"priority"`
	Target   string `json:"target"`
	Icon     string `json:"icon"`
	ParentID int32  `json:"parentId"`
	Team     string `json:"team"`
}
