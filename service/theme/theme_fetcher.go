package theme

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
)

type ThemeFetcher interface {
	FetchTheme(ctx context.Context, file interface{}) (*dto.ThemeProperty, error)
}
