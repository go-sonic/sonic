package extension

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

type photoExtension struct {
	PhotoService  service.PhotoService
	OptionService service.OptionService
	Template      *template.Template
}

func RegisterPhotoFunc(template *template.Template, photoService service.PhotoService, optionService service.OptionService) {
	p := &photoExtension{
		Template:      template,
		PhotoService:  photoService,
		OptionService: optionService,
	}
	p.addGetPhotoCount()
	p.addListByTeam()
	p.addListPhotos()
	p.addListPhotosGroupByTeam()
}

func (p *photoExtension) addListPhotos() {
	listPhotos := func() ([]*dto.Photo, error) {
		ctx := context.Background()
		photos, err := p.PhotoService.List(ctx, &param.Sort{
			Fields: []string{"createTime,desc"},
		})
		if err != nil {
			return nil, err
		}
		return p.PhotoService.ConvertToDTOs(ctx, photos), nil
	}
	p.Template.AddFunc("listPhotos", listPhotos)
}

func (p *photoExtension) addListPhotosGroupByTeam() {
	listPhotosGroupByTeam := func() (map[string][]*dto.Photo, error) {
		ctx := context.Background()
		photos, err := p.PhotoService.List(ctx, &param.Sort{
			Fields: []string{"createTime,desc"},
		})
		if err != nil {
			return nil, err
		}
		teamPhotos := make(map[string][]*dto.Photo)
		for _, photo := range photos {
			teamPhotos[photo.Team] = append(teamPhotos[photo.Team], p.PhotoService.ConvertToDTO(ctx, photo))
		}
		return teamPhotos, nil
	}
	p.Template.AddFunc("listTeamPhotos", listPhotosGroupByTeam)
}

func (p *photoExtension) addListByTeam() {
	listPhotoByTeam := func(team string) ([]*dto.Photo, error) {
		ctx := context.Background()
		photos, err := p.PhotoService.ListByTeam(ctx, team, &param.Sort{
			Fields: []string{"createTime,desc"},
		})
		return p.PhotoService.ConvertToDTOs(ctx, photos), err
	}
	p.Template.AddFunc("listPhotoByTeam", listPhotoByTeam)
}

func (p *photoExtension) addGetPhotoCount() {
	getPhotoCount := func() (int64, error) {
		ctx := context.Background()
		return p.PhotoService.GetPhotoCount(ctx)
	}
	p.Template.AddFunc("getPhotoCount", getPhotoCount)
}
