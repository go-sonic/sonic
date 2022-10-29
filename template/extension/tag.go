package extension

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

type tagExtension struct {
	PostTagService service.PostTagService
	TagService     service.TagService
	Template       *template.Template
}

func RegisterTagFunc(template *template.Template, postTagService service.PostTagService, tagService service.TagService) {
	te := tagExtension{
		PostTagService: postTagService,
		TagService:     tagService,
		Template:       template,
	}
	te.addGetAllTag()
	te.addGetTagByPostID()
	te.addGetTagByPostID()
	te.addGetTagCount()
}

func (te *tagExtension) addGetAllTag() {
	getAllTag := func() ([]*dto.TagWithPostCount, error) {
		return te.PostTagService.ListAllTagWithPostCount(context.Background(), &param.Sort{
			Fields: []string{"createTime,desc"},
		})
	}
	te.Template.AddFunc("getAllTag", getAllTag)
}

func (te *tagExtension) addGetTagByPostID() {
	getTagByPostID := func(postID int32) ([]*dto.Tag, error) {
		ctx := context.Background()

		tags, err := te.PostTagService.ListTagByPostID(ctx, postID)
		if err != nil {
			return nil, err
		}
		return te.TagService.ConvertToDTOs(ctx, tags)
	}
	te.Template.AddFunc("getTagByPostID", getTagByPostID)
}

func (te *tagExtension) addGetTagCount() {
	getTagCount := func() (int64, error) {
		ctx := context.Background()
		return te.TagService.CountAllTag(ctx)
	}
	te.Template.AddFunc("getTagCount", getTagCount)
}
