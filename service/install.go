package service

import (
	"context"

	"github.com/go-sonic/sonic/model/param"
)

type InstallService interface {
	InstallBlog(ctx context.Context, installParam param.Install) error
}
