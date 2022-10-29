package vo

import "github.com/go-sonic/sonic/model/dto"

type Menu struct {
	dto.Menu
	Children []*Menu `json:"children"`
}
