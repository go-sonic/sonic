package extension

import (
	"context"
	"math/rand"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

type linkExtension struct {
	LinkService service.LinkService
	Template    *template.Template
}

func RegisterLinkFunc(template *template.Template, linkService service.LinkService) {
	l := &linkExtension{
		LinkService: linkService,
		Template:    template,
	}
	l.addListLinks()
	l.addGetLinksCount()
	l.addListLinksGroupByTeam()
	l.addListLinksRandom()
}

func (l *linkExtension) addListLinks() {
	listLinks := func() ([]*dto.Link, error) {
		ctx := context.Background()
		links, err := l.LinkService.List(ctx, nil)
		if err != nil {
			return nil, err
		}
		return l.LinkService.ConvertToDTOs(ctx, links), nil
	}
	l.Template.AddFunc("listLinks", listLinks)
}

func (l *linkExtension) addListLinksRandom() {
	listLinksRandom := func() ([]*dto.Link, error) {
		ctx := context.Background()
		links, err := l.LinkService.List(ctx, nil)
		if err != nil {
			return nil, err
		}
		rand.Shuffle(len(links), func(i, j int) {
			links[i], links[j] = links[j], links[i]
		})
		return l.LinkService.ConvertToDTOs(ctx, links), nil
	}
	l.Template.AddFunc("listLinksRandom", listLinksRandom)
}

func (l *linkExtension) addListLinksGroupByTeam() {
	listLinksGroupByTeam := func() (map[string][]*dto.Link, error) {
		ctx := context.Background()
		links, err := l.LinkService.List(ctx, nil)
		if err != nil {
			return nil, err
		}
		linkDTOs := l.LinkService.ConvertToDTOs(ctx, links)
		teamLinkMap := make(map[string][]*dto.Link)
		for _, link := range linkDTOs {
			teamLinkMap[link.Team] = append(teamLinkMap[link.Team], link)
		}
		return teamLinkMap, nil
	}
	l.Template.AddFunc("listLinksGroupByTeam", listLinksGroupByTeam)
}

func (l *linkExtension) addGetLinksCount() {
	getLinksCount := func() (int64, error) {
		ctx := context.Background()
		return l.LinkService.Count(ctx)
	}
	l.Template.AddFunc("getLinksCount", getLinksCount)
}
