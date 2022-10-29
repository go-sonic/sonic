package vo

import "github.com/go-sonic/sonic/model/dto"

type CategoryVO struct {
	dto.CategoryDTO
	Children []*CategoryVO `json:"children"`
}
