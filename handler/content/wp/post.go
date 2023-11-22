package wp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	sonicconst "github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto/wp"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
	"strconv"
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

func (handler *PostHandler) List(ctx *gin.Context) (interface{}, error) {
	var wpPostQuery param.WpPostQuery
	if err := ctx.ShouldBind(&wpPostQuery); err != nil {
		return nil, util.WrapJsonBindErr(err)
	}

	var postQuery param.PostQuery
	postQuery.PageSize = wpPostQuery.Page
	postQuery.PageNum = wpPostQuery.PerPage
	entities, _, err := handler.PostService.Page(ctx, postQuery)
	if err != nil {
		return nil, err
	}

	wpPostList := make([]*wp.PostDTO, 0, len(entities))
	for _, postEntity := range entities {
		wpPostList = append(wpPostList, convertToWpPost(postEntity))
	}

	return wpPostList, nil
}

func (handler *PostHandler) Create(ctx *gin.Context) (interface{}, error) {
	postParam, err := parsePostParam(ctx)
	if err != nil {
		return nil, util.WrapJsonBindErr(err)
	}

	create, err := handler.PostService.Create(ctx, postParam)
	if err != nil {
		return nil, err
	}
	return convertToWpPost(create), nil
}

func (handler *PostHandler) Update(ctx *gin.Context) (interface{}, error) {
	postParam, err := parsePostParam(ctx)
	if err != nil {
		return nil, util.WrapJsonBindErr(err)
	}

	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}

	postDetail, err := handler.PostService.Update(ctx, int32(postID), postParam)
	if err != nil {
		return nil, err
	}

	return convertToWpPost(postDetail), nil
}

func (handler *PostHandler) Delete(ctx *gin.Context) (interface{}, error) {
	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}

	postEntity, err := handler.PostService.GetByPostID(ctx, int32(postID))
	if err != nil {
		return nil, err
	}

	if err = handler.PostService.Delete(ctx, int32(postID)); err != nil {
		return nil, err
	}

	return convertToWpPost(postEntity), nil
}

func parsePostParam(ctx *gin.Context) (*param.Post, error) {
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

	return convertToPostParam(&wpPost)

}

func convertToPostParam(wpPost *param.WpPost) (*param.Post, error) {
	var paramPostStatus = sonicconst.PostStatusPublished
	if strings.ToLower(wpPost.Status) == "draft" {
		paramPostStatus = sonicconst.PostStatusDraft
	}

	editorType := sonicconst.EditorTypeRichText
	disallowComment := false
	if strings.ToLower(wpPost.CommentStatus) == "closed" {
		disallowComment = true
	}

	var createTime *int64

	if strings.TrimSpace(wpPost.Date) != "" {
		datetime, err := time.Parse(sonicconst.LocalDateTimeFormat, wpPost.Date)
		if err != nil {
			return nil, err
		}
		dateTimeMills := datetime.UnixMilli()
		createTime = &dateTimeMills
	}

	return &param.Post{
		Title:           wpPost.Title,
		Status:          paramPostStatus,
		Slug:            wpPost.Slug,
		EditorType:      &editorType,
		OriginalContent: wpPost.Content,
		Summary:         "",
		Thumbnail:       "",
		DisallowComment: disallowComment,
		Password:        wpPost.Password,
		Template:        "",
		TopPriority:     0,
		CreateTime:      createTime,
		MetaKeywords:    "",
		MetaDescription: "",
		TagIDs:          make([]int32, 0),
		CategoryIDs:     make([]int32, 0),
		MetaParam:       nil,
		Content:         wpPost.Content,
		EditTime:        nil,
		UpdateTime:      nil,
	}, nil
}

func convertToWpPost(postEntity *entity.Post) *wp.PostDTO {
	timeFormat := time.RFC3339
	var wpStatus = "publish"
	var wpCommentStatus = "open"
	var wpContent = make(map[string]interface{})

	if postEntity.Status == sonicconst.PostStatusDraft {
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
		Modified:          "",
		ModifiedGmt:       "",
		Slug:              "",
		Status:            wpStatus,
		Type:              "post",
		Password:          "standard",
		PermalinkTemplate: "",
		GeneratedSlug:     "",
		Title:             postEntity.Title,
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

	if postEntity.UpdateTime != nil {
		postDTO.Modified = postEntity.UpdateTime.Format(timeFormat)
		postDTO.ModifiedGmt = postEntity.UpdateTime.UTC().Format(timeFormat)
	}
	return postDTO
}
