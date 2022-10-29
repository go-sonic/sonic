package dto

type IndependentSheet struct {
	ID        int32  `json:"id"`
	Title     string `json:"title"`
	FullPath  string `json:"fullPath"`
	RouteName string `json:"routeName"`
	Available bool   `json:"available"`
}
