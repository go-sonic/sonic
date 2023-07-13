package theme

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/fx"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type multipartZipThemeFetcherImpl struct {
	fx.Out
	PropertyScanner PropertyScanner
}

func NewMultipartZipThemeFetcher(propertyScanner PropertyScanner) ThemeFetcher {
	return &multipartZipThemeFetcherImpl{
		PropertyScanner: propertyScanner,
	}
}

func (m *multipartZipThemeFetcherImpl) FetchTheme(ctx context.Context, file interface{}) (*dto.ThemeProperty, error) {
	themeFileHeader, ok := file.(*multipart.FileHeader)
	if !ok {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("not support")
	}

	tempDir := os.TempDir()
	srcThemeFile, err := themeFileHeader.Open()
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("upload theme file error")
	}
	defer srcThemeFile.Close()

	fileName := themeFileHeader.Filename
	if !strings.HasSuffix(fileName, ".zip") {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("not zip file")
	}

	diskFilePath := filepath.Join(tempDir, fileName)
	if util.FileIsExisted(diskFilePath) {
		err = os.Remove(diskFilePath)
		if err != nil {
			return nil, xerr.WithStatus(err, xerr.StatusInternalServerError)
		}
	}

	diskFile, err := os.OpenFile(diskFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o444)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("create file error")
	}

	defer diskFile.Close()

	_, err = io.Copy(diskFile, srcThemeFile)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("save file error")
	}
	_, err = util.Unzip(filepath.Join(tempDir, fileName), filepath.Join(tempDir, strings.TrimSuffix(fileName, ".zip")))
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("unzip file error")
	}
	dest := filepath.Join(tempDir, strings.TrimSuffix(fileName, ".zip"))
	themeProperty, err := m.PropertyScanner.ReadThemeProperty(ctx, dest)
	if err == nil && themeProperty != nil {
		return themeProperty, nil
	}
	dirEntrys, err := os.ReadDir(dest)
	for _, dirEntry := range dirEntrys {
		if !dirEntry.IsDir() {
			continue
		}
		themeProperty, err = m.PropertyScanner.ReadThemeProperty(ctx, filepath.Join(dest, dirEntry.Name()))
		if err == nil && themeProperty != nil {
			return themeProperty, nil
		}
	}
	return nil, err
}
