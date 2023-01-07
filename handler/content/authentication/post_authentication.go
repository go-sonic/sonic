package authentication

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type PostAuthentication struct {
	OptionService       service.OptionService
	PostService         service.PostService
	PostCategoryService service.PostCategoryService
	CategoryService     service.CategoryService
}

func NewPostAuthentication(
	optionService service.OptionService,
	postService service.PostService,
	categoryService service.CategoryService,
	postCategoryService service.PostCategoryService,
) *PostAuthentication {
	return &PostAuthentication{
		CategoryService:     categoryService,
		OptionService:       optionService,
		PostService:         postService,
		PostCategoryService: postCategoryService,
	}
}

func (p *PostAuthentication) Authenticate(ctx context.Context, token string, id int32, password string) (string, error) {
	post, err := p.PostService.GetByPostID(ctx, id)
	if err != nil {
		return "", err
	}
	if post.Password != "" {
		if post.Password == password {
			return p.doAuthenticate(ctx, token, id)
		} else {
			return "", xerr.WithMsg(nil, "密码不正确").WithStatus(http.StatusUnauthorized)
		}
	}
	postCategories, err := p.PostCategoryService.ListCategoryByPostID(ctx, id)
	if err != nil {
		return "", err
	}

	for _, category := range postCategories {
		if category.Password == password {
			return p.doAuthenticate(ctx, token, id)
		}
	}

	allCategories, err := p.CategoryService.ListAll(ctx, nil)
	if err != nil {
		return "", err
	}
	categoryMap := make(map[int32]*entity.Category)
	for _, category := range allCategories {
		categoryMap[category.ID] = category
	}

	for _, postCategory := range postCategories {
		parentID := postCategory.ParentID
		for {
			parentCategory, ok := categoryMap[parentID]
			if !ok || parentCategory == nil {
				break
			}
			switch parentCategory.Password {
			case "":
				parentID = parentCategory.ParentID
			case password:
				return p.doAuthenticate(ctx, token, id)
			default:
				break
			}
		}
	}
	return "", xerr.WithMsg(nil, "密码不正确").WithStatus(http.StatusUnauthorized)
}

func (p *PostAuthentication) IsAuthenticated(ctx context.Context, tokenStr string, id int32) (bool, error) {
	if tokenStr == "" {
		return false, nil
	}

	secret, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.JWTSecret, "")
	if err != nil {
		return false, err
	}
	if secret.(string) == "" {
		return false, xerr.WithMsg(nil, "jwt secret is nil").WithStatus(xerr.StatusInternalServerError)
	}

	post, err := p.PostService.GetByPostID(ctx, id)
	if err != nil {
		return false, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret.(string)), nil
	})
	if err != nil {
		return false, err
	}
	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid || claims == nil {
		return false, nil
	}

	for _, postID := range claims.PostIDs {
		if postID == id {
			return true, nil
		}
	}
	if post.Password != "" {
		return false, nil
	}
	categories, err := p.PostCategoryService.ListCategoryByPostID(ctx, id)
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}
	for _, categoryID := range claims.CategoryIDs {
		for _, category := range categories {
			if category.Type == consts.CategoryTypeNormal {
				continue
			}
			if category.ID == categoryID {
				return true, nil
			}
		}
	}
	return false, nil
}

func (p *PostAuthentication) doAuthenticate(ctx context.Context, tokenStr string, id int32) (string, error) {
	secret, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.JWTSecret, "")
	if err != nil {
		return "", err
	}
	if secret.(string) == "" {
		return "", xerr.WithMsg(nil, "jwt secret is nil").WithStatus(xerr.StatusInternalServerError)
	}
	var claims *customClaims
	if tokenStr != "" {
		token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret.(string)), nil
		})
		if err == nil {
			if c, ok := token.Claims.(*customClaims); ok && token.Valid {
				claims = c
			}
		}
	}
	if claims == nil {
		claims = &customClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}
	}
	claims.PostIDs = append(claims.PostIDs, id)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret.(string)))
	if err != nil {
		return "", xerr.WithStatus(err, xerr.StatusInternalServerError)
	}
	return ss, nil
}
