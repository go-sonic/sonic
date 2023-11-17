package wp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto/wp"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"strings"
	"time"
)

type PostHandler struct {
	PostService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		PostService: postService,
	}
}

func (handler *PostHandler) Create(ctx *gin.Context) (interface{}, error) {
	var wpPost param.WpPost
	err := ctx.ShouldBindJSON(&wpPost)
	if err != nil {
		return nil, util.WrapJsonBindErr(err)
	}

	bytes, err := json.Marshal(wpPost)
	if err != nil {
		return nil, err
	}
	log.CtxInfo(ctx, "wpPost: "+string(bytes))

	postParam := convertToPostParam(&wpPost)

	create, err := handler.PostService.Create(ctx, postParam)
	if err != nil {
		return nil, err
	}

	return convertToWpPost(create), nil
}

func convertToPostParam(wpPost *param.WpPost) *param.Post {
	var paramPostStatus = consts.PostStatusPublished
	if strings.ToLower(wpPost.Content) == "draft" {
		paramPostStatus = consts.PostStatusDraft
	}

	createTime := time.Now().Unix()

	return &param.Post{
		Title:           wpPost.Title,
		Status:          paramPostStatus,
		Slug:            wpPost.Slug,
		EditorType:      nil,
		OriginalContent: wpPost.Content,
		Summary:         "",
		Thumbnail:       "",
		DisallowComment: false,
		Password:        wpPost.Password,
		Template:        "",
		TopPriority:     0,
		CreateTime:      &createTime,
		MetaKeywords:    "",
		MetaDescription: "",
		TagIDs:          make([]int32, 0),
		CategoryIDs:     make([]int32, 0),
		MetaParam:       nil,
		Content:         wpPost.Content,
		EditTime:        nil,
		UpdateTime:      nil,
	}
}

func convertToWpPost(postEntity *entity.Post) *wp.PostDTO {
	timeFormat := time.RFC3339
	var wpStatus = "publish"
	var wpCommentStatus = "open"
	var wpContent = make(map[string]interface{})

	if postEntity.Status == consts.PostStatusDraft {
		wpStatus = "draft"
	}

	if postEntity.DisallowComment {
		wpCommentStatus = "close"
	}

	wpContent["rendered"] = postEntity.OriginalContent
	wpContent["protected"] = false

	var postDTO = &wp.PostDTO{
		Date:              postEntity.CreateTime.Format(timeFormat),
		DateGmt:           postEntity.CreateTime.UTC().Format(timeFormat),
		Guid:              nil,
		Id:                postEntity.ID,
		Link:              "",
		Modified:          postEntity.UpdateTime.Format(timeFormat),
		ModifiedGmt:       postEntity.UpdateTime.UTC().Format(timeFormat),
		Slug:              "",
		Status:            wpStatus,
		Type:              "post",
		Password:          "standard",
		PermalinkTemplate: "",
		GeneratedSlug:     "",
		Title:             "",
		Content:           wpContent,
		Author:            0,
		Excerpt:           nil,
		FeaturedMedia:     0,
		CommentStatus:     wpCommentStatus,
		PingStatus:        "open",
		Format:            "standard",
		Meta:              nil,
		Sticky:            false,
		Template:          "",
		Categories:        make([]int32, 0),
		Tags:              make([]int32, 0),
	}
	return postDTO
}
