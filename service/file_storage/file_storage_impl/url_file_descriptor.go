package file_storage_impl

import (
	"net/url"
	"path/filepath"

	"github.com/go-sonic/sonic/util/xerr"
)

type urlFileDescriptor struct {
	OriginalName    string
	Name            string
	Extension       string
	BasePath        string
	SubPath         string
	AutomaticRename bool
	ShouldRename    func(relativePath string) (bool, error)
	Suffix          string
}

func newURLFileDescriptor(opts ...urlOption) (fileDescriptor, error) {
	fd := &urlFileDescriptor{}
	for _, opt := range opts {
		opt(fd)
	}
	_, err := url.Parse(fd.BasePath)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	_, err = url.Parse(fd.SubPath)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}

	fd.OriginalName = filepath.Clean(fd.OriginalName)

	ext := filepath.Ext(fd.OriginalName)
	if ext != "" {
		// remove dot
		fd.Extension = ext[1:]
	}
	originalName := filepath.Base(fd.OriginalName)
	for i := len(originalName) - 1; i >= 0; i-- {
		if originalName[i] == '.' {
			fd.Name = originalName[0:i]
			break
		}
	}
	err = rename(fd)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

func (f *urlFileDescriptor) getFullName() string {
	if f.Extension == "" {
		return f.Name + f.Suffix
	}
	return f.Name + f.Suffix + "." + f.Extension
}

func (f *urlFileDescriptor) getFullDirPath() string {
	panic("not support")
}

func (f *urlFileDescriptor) getFullPath() string {
	fullPath, _ := url.JoinPath(f.BasePath, f.SubPath, f.getFullName())
	return fullPath
}

func (f *urlFileDescriptor) getRelativePath() string {
	relativePath, _ := url.JoinPath(f.SubPath, f.getFullName())
	return relativePath
}

func (f *urlFileDescriptor) getExtensionName() string {
	return f.Extension
}

func (f *urlFileDescriptor) getFileName() string {
	return f.Name
}

func (f *urlFileDescriptor) setFileName(name string) {
	f.Name = name
}

func (f *urlFileDescriptor) getShouldRename() shouldRename {
	return f.ShouldRename
}

type urlOption func(f *urlFileDescriptor)

func withBaseURL(baseURL string) urlOption {
	return func(f *urlFileDescriptor) {
		f.BasePath = baseURL
	}
}

func withSubURLPath(subURL string) urlOption {
	return func(f *urlFileDescriptor) {
		f.SubPath = subURL
	}
}

func withOriginalNameURLOption(originalName string) urlOption {
	return func(f *urlFileDescriptor) {
		f.OriginalName = originalName
	}
}

func withShouldRenameURLOption(fn func(relativePath string) (bool, error)) urlOption {
	return func(f *urlFileDescriptor) {
		f.ShouldRename = fn
	}
}
