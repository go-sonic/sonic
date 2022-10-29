package vo

import "github.com/go-sonic/sonic/model/dto"

type LinkTeamVO struct {
	Team  string
	Links []*dto.Link
}
