package service

import (
	"context"
	"mime/multipart"

	"github.com/go-sonic/sonic/model/dto"
)

type BackupService interface {
	// GetBackup  Get backup data by backup file name.
	GetBackup(ctx context.Context, filename string, backupType BackupType) (*dto.BackupDTO, error)
	// BackupWholeSite Zips work directory
	BackupWholeSite(ctx context.Context, toBackupItems []string) (*dto.BackupDTO, error)
	// ListFiles list all files under path
	ListFiles(ctx context.Context, path string, backupType BackupType) ([]*dto.BackupDTO, error)
	// GetBackupFilePath get filepath and check if the file exist
	GetBackupFilePath(ctx context.Context, path string, filename string) (string, error)
	// DeleteFile delete file
	DeleteFile(ctx context.Context, path string, filename string) error
	// ExportData export database data to json file
	ExportData(ctx context.Context) (*dto.BackupDTO, error)
	// ImportMarkdown import markdown file as post
	ImportMarkdown(ctx context.Context, fileHeader *multipart.FileHeader) error
	// ExportMarkdown export posts to markdown files
	ExportMarkdown(ctx context.Context, needFrontMatter bool) (*dto.BackupDTO, error)
	ListToBackupItems(ctx context.Context) ([]string, error)
}

type BackupType string

const (
	WholeSite BackupType = "/api/admin/backups/work-dir"
	JsonData  BackupType = "/api/admin/backups/data"
	Markdown  BackupType = "/api/admin/backups/markdown/export"
)
