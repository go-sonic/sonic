package template

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func (t *Template) Watch() {
	for {
		select {
		case event, ok := <-t.watcher.Events:
			if !ok {
				return
			}
			if filepath.Ext(event.Name) != ".tmpl" {
				continue
			}

			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				t.logger.Info("Write file:", zap.String("file", event.Name))
			case event.Op&fsnotify.Create == fsnotify.Create:
				t.logger.Info("Create file:", zap.String("file", event.Name))
			default:
				continue
			}

			err := t.Reload([]string{event.Name})
			if err != nil {
				t.logger.Error("reload template error", zap.Error(err))
			}
		case err, ok := <-t.watcher.Errors:
			if !ok {
				return
			}
			t.logger.Error("file watcher error:", zap.Error(err))
		}
	}
}
