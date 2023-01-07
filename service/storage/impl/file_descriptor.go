package filestorageimpl

import (
	"path/filepath"
	"strconv"
	"time"
)

type fileDescriptor interface {
	getFullDirPath() string
	getFullPath() string
	getRelativePath() string
	getFileName() string
	setFileName(string)
	getExtensionName() string
	getShouldRename() shouldRename
}

type shouldRename = func(relativePath string) (bool, error)

type localFileDescriptor struct {
	OriginalName    string
	Name            string
	Extension       string
	BasePath        string
	SubPath         string
	AutomaticRename bool
	ShouldRename    shouldRename
	Suffix          string
}

func newLocalFileDescriptor(opts ...fdOption) (fileDescriptor, error) {
	fd := &localFileDescriptor{}
	for _, opt := range opts {
		opt(fd)
	}
	fd.BasePath = filepath.Clean(fd.BasePath)
	fd.SubPath = filepath.Clean(fd.SubPath)
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
	err := rename(fd)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

func rename(f fileDescriptor) error {
	if f.getShouldRename() == nil {
		return nil
	}
	should, err := f.getShouldRename()(f.getRelativePath())
	if err != nil {
		return err
	}
	if !should {
		return nil
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	f.setFileName(f.getFileName() + "-" + timestamp)
	return nil
}

func (f *localFileDescriptor) getFullName() string {
	if f.Extension == "" {
		return f.Name + f.Suffix
	}
	return f.Name + f.Suffix + "." + f.Extension
}

func (f *localFileDescriptor) setFileName(name string) {
	f.Name = name
}

func (f *localFileDescriptor) getFullDirPath() string {
	return filepath.Join(f.BasePath, f.SubPath)
}

func (f *localFileDescriptor) getFullPath() string {
	return filepath.Join(f.BasePath, f.SubPath, f.getFullName())
}

func (f *localFileDescriptor) getRelativePath() string {
	return filepath.Join(f.SubPath, f.getFullName())
}

func (f *localFileDescriptor) getExtensionName() string {
	return f.Extension
}

func (f *localFileDescriptor) getFileName() string {
	return f.Name
}

func (f *localFileDescriptor) getShouldRename() shouldRename {
	return f.ShouldRename
}

type fdOption func(f *localFileDescriptor)

func withOriginalName(originalName string) fdOption {
	return func(f *localFileDescriptor) {
		f.OriginalName = originalName
	}
}

func withBasePath(basePath string) fdOption {
	return func(f *localFileDescriptor) {
		f.BasePath = basePath
	}
}

func withSubPath(subPath string) fdOption {
	return func(f *localFileDescriptor) {
		f.SubPath = subPath
	}
}

func withShouldRename(fn func(relativePath string) (bool, error)) fdOption {
	return func(f *localFileDescriptor) {
		f.ShouldRename = fn
	}
}

func withSuffix(suffix string) fdOption {
	return func(f *localFileDescriptor) {
		f.Suffix = suffix
	}
}
