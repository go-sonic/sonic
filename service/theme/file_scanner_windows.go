//go:build windows

package theme

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/util/xerr"
)

type FileScanner interface {
	ListThemeFiles(ctx context.Context, themePath string) ([]*dto.ThemeFile, error)
}

func NewFileScanner() FileScanner {
	return &fileScannerImpl{}
}

type fileScannerImpl struct{}

func (f *fileScannerImpl) ListThemeFiles(ctx context.Context, themePath string) ([]*dto.ThemeFile, error) {
	fileMap := make(map[string]*dto.ThemeFile)
	root := &dto.ThemeFile{}
	fileMap[themePath] = root

	err := filepath.Walk(themePath, func(path string, info fs.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return err
		} else if !os.IsPermission(err) && err != nil {
			return err
		}
		themeFile := &dto.ThemeFile{
			Name:     info.Name(),
			IsFile:   !info.IsDir(),
			Path:     path,
			Editable: true,
		}
		parentDir, ok := fileMap[filepath.Dir(path)]
		if !ok {
			return nil
		}
		parentDir.Node = append(parentDir.Node, themeFile)

		fileMap[path] = themeFile
		return nil
	})
	if err != nil {
		return nil, xerr.NoType.Wrap(err)
	}
	return root.Node, nil
}
