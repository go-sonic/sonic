package config

import (
	"os"
	"path/filepath"
)

var (
	TempDir           = os.TempDir()
	BackupDir         = filepath.Join(TempDir, "sonic-backup") + string(os.PathSeparator)
	BackupMarkdownDir = filepath.Join(TempDir, "sonic-backup-markdown") + string(os.PathSeparator)
	DataExportDir     = filepath.Join(TempDir, "sonic-data-export") + string(os.PathSeparator)
	ResourcesDir, _   = filepath.Abs("./resources")
)
