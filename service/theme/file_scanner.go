//go:build linux || darwin || freebsd || netbsd || openbsd || dragonfly

package theme

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

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

	euid, egid := os.Geteuid(), os.Getegid()
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
			Editable: f.isFileWritable(info.Sys().(*syscall.Stat_t), info.Mode(), euid, egid) && f.isFileReadable(info.Sys().(*syscall.Stat_t), info.Mode(), euid, egid),
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

func (f *fileScannerImpl) isFileWritable(t *syscall.Stat_t, fileMode fs.FileMode, euid, egid int) bool {
	if t == nil {
		return true
	}
	if t.Uid == uint32(euid) {
		return (uint32(fileMode) & uint32(0o200)) != 0
	}
	if t.Gid == uint32(egid) {
		return (uint32(fileMode) & uint32(0o020)) != 0
	} else {
		return (uint32(fileMode) & uint32(0o002)) != 0
	}
}

func (f *fileScannerImpl) isFileReadable(t *syscall.Stat_t, fileMode fs.FileMode, euid, egid int) bool {
	if t == nil {
		return true
	}
	if t.Uid == uint32(euid) {
		return (uint32(fileMode) & uint32(0o400)) != 0
	}
	if t.Gid == uint32(egid) {
		return (uint32(fileMode) & uint32(0o040)) != 0
	} else {
		return (uint32(fileMode) & uint32(0o004)) != 0
	}
}
