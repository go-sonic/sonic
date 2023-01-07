package authentication

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type CategoryAuthentication struct {
	OptionService   service.OptionService
	CategoryService service.CategoryService
}

func NewCategoryAuthentication(
	optionService service.OptionService,
	categoryService service.CategoryService,
) *CategoryAuthentication {
	return &CategoryAuthentication{
		OptionService:   optionService,
		CategoryService: categoryService,
	}
}

// Authenticate implements ContentAuthentication
func (c *CategoryAuthentication) Authenticate(ctx context.Context, token string, id int32, password string) (string, error) {
	category, err := c.CategoryService.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	if category.Password == "" && category.ParentID == 0 {
		return "", nil
	}
	if category.Password == "" {
		categories, err := c.CategoryService.ListAll(ctx, nil)
		if err != nil {
			return "", err
		}
		categoryMap := make(map[int32]*entity.Category)
		for _, category := range categories {
			categoryMap[category.ID] = category
		}
		parentID := category.ParentID
		parentIDs := make([]int32, 0)
		for {
			parentCategory, ok := categoryMap[parentID]
			if !ok || parentCategory == nil {
				return "", nil
			}
			switch parentCategory.Password {
			case "":
				parentID = parentCategory.ParentID
				parentIDs = append(parentIDs, parentID)
			case password:
				return c.doAuthenticate(ctx, token, parentIDs...)
			default:
				return "", xerr.WithMsg(nil, "密码不正确").WithStatus(http.StatusUnauthorized)
			}
		}
	} else if category.Password == password {
		return c.doAuthenticate(ctx, token, id)
	}
	return "", xerr.WithMsg(nil, "密码不正确").WithStatus(http.StatusUnauthorized)
}

func (c *CategoryAuthentication) IsAuthenticated(ctx context.Context, tokenStr string, id int32) (bool, error) {
	if tokenStr == "" {
		return false, nil
	}
	secret, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.JWTSecret, "")
	if err != nil {
		return false, err
	}
	if secret.(string) == "" {
		return false, xerr.WithMsg(nil, "jwt secret is nil").WithStatus(xerr.StatusInternalServerError)
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

	for _, categoryID := range claims.CategoryIDs {
		if categoryID == id {
			return true, nil
		}
	}

	return false, nil
}

func (c *CategoryAuthentication) doAuthenticate(ctx context.Context, tokenStr string, id ...int32) (string, error) {
	secret, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.JWTSecret, "")
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
	claims.CategoryIDs = append(claims.CategoryIDs, id...)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret.(string)))
	if err != nil {
		return "", xerr.WithStatus(err, xerr.StatusInternalServerError)
	}
	return ss, nil
}
