package impl

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type backupServiceImpl struct {
	Config              *config.Config
	OptionService       service.OptionService
	OneTimeTokenService service.OneTimeTokenService
}

func NewBackUpService(config *config.Config, optionService service.OptionService, oneTimeTokenService service.OneTimeTokenService) service.BackupService {
	return &backupServiceImpl{
		Config:              config,
		OptionService:       optionService,
		OneTimeTokenService: oneTimeTokenService,
	}
}

func (b *backupServiceImpl) GetBackup(ctx context.Context, filepath string, backupType service.BackupType) (*dto.BackupDTO, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return nil, xerr.NoType.Wrap(err).WithMsg("file not exist")
	} else if err != nil {
		return nil, xerr.NoType.Wrap(err)
	}
	return b.buildBackupDTO(ctx, string(backupType), filepath)
}

func (b *backupServiceImpl) BackupWholeSite(ctx context.Context, toBackupItems []string) (*dto.BackupDTO, error) {
	backupFilename := consts.SonicBackupPrefix + time.Now().Format("2006-01-02-15-04-05") + util.GenUUIDWithOutDash() + ".zip"
	backupFilePath := config.BackupDir

	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		err = os.MkdirAll(backupFilePath, os.ModePerm)
		if err != nil {
			return nil, xerr.NoType.Wrap(err).WithMsg("create dir err")
		}
	} else if err != nil {
		return nil, xerr.NoType.Wrap(err).WithMsg("get fileInfo")
	}

	backupFile := filepath.Join(backupFilePath, backupFilename)

	toBackupPaths := []string{}
	for _, toBackupItem := range toBackupItems {
		toBackupPath := filepath.Join(b.Config.Sonic.WorkDir, toBackupItem)
		toBackupPaths = append(toBackupPaths, toBackupPath)
	}

	err := util.ZipFile(backupFile, toBackupPaths...)
	if err != nil {
		return nil, err
	}
	return b.buildBackupDTO(ctx, string(service.WholeSite), backupFile)
}

func (b *backupServiceImpl) ListFiles(ctx context.Context, path string, backupType service.BackupType) ([]*dto.BackupDTO, error) {
	backups := make([]*dto.BackupDTO, 0)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return backups, nil
	} else if err != nil {
		return nil, xerr.NoType.Wrap(err).WithMsg("Failed to fetch backups")
	}
	prefix := ""
	if backupType == service.WholeSite {
		prefix = consts.SonicBackupPrefix
	} else if backupType == service.JsonData {
		prefix = consts.SonicDataExportPrefix
	} else if backupType == service.Markdown {
		prefix = consts.SonicBackupMarkdownPrefix
	}
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), prefix) {
			backupDTO, err := b.buildBackupDTO(ctx, string(backupType), path)
			if err != nil {
				return err
			}
			backups = append(backups, backupDTO)
		}
		return nil
	})
	if err != nil {
		return nil, xerr.NoType.Wrap(err).WithMsg("Failed to fetch backups")
	}
	return backups, nil
}

func (b *backupServiceImpl) GetBackupFilePath(ctx context.Context, path string, filename string) (string, error) {
	backupFilePath := filepath.Join(path, filename)
	_, err := os.Stat(backupFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", xerr.NoType.Wrap(err).WithStatus(xerr.StatusBadRequest).WithMsg("file not exist")
		}
		return "", xerr.NoType.Wrap(err).WithStatus(xerr.StatusInternalServerError)
	}
	return backupFilePath, nil
}

func (b *backupServiceImpl) DeleteFile(ctx context.Context, path string, filename string) error {
	backupFilePath := filepath.Join(path, filename)
	err := os.Remove(backupFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return xerr.NoType.Wrap(err).WithMsg("file not exist").WithStatus(xerr.StatusBadRequest)
		}
		return xerr.NoType.Wrap(err).WithMsg("Failed to delete file")
	}
	return nil
}

func (b *backupServiceImpl) ImportMarkdown(ctx context.Context, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return xerr.NoType.Wrap(err).WithMsg("upload file error")
	}
	bContent, err := io.ReadAll(file)
	if err != nil {
		return xerr.NoType.Wrap(err).WithMsg("read file error")
	}
	content := string(bContent)
	log.Info(content)
	// TODO 导入markdown
	return nil
}

func (b *backupServiceImpl) ExportData(ctx context.Context) (*dto.BackupDTO, error) {
	data := make(map[string]interface{})
	data["version"] = consts.SonicVersion
	data["export_date"] = time.Now().Format("2006-01-02 15:04:05")
	err := fillData(data, "attachments", dal.GetQueryByCtx(ctx).Attachment.WithContext(ctx).Find, nil)
	err = fillData(data, "category", dal.GetQueryByCtx(ctx).Category.WithContext(ctx).Find, err)
	err = fillData(data, "comment", dal.GetQueryByCtx(ctx).Comment.WithContext(ctx).Find, err)
	err = fillData(data, "comment_black", dal.GetQueryByCtx(ctx).CommentBlack.WithContext(ctx).Find, err)
	err = fillData(data, "journal", dal.GetQueryByCtx(ctx).Journal.WithContext(ctx).Find, err)
	err = fillData(data, "link", dal.GetQueryByCtx(ctx).Link.WithContext(ctx).Find, err)
	err = fillData(data, "log", dal.GetQueryByCtx(ctx).Log.WithContext(ctx).Find, err)
	err = fillData(data, "menu", dal.GetQueryByCtx(ctx).Menu.WithContext(ctx).Find, err)
	err = fillData(data, "meta", dal.GetQueryByCtx(ctx).Meta.WithContext(ctx).Find, err)
	err = fillData(data, "option", dal.GetQueryByCtx(ctx).Option.WithContext(ctx).Find, err)
	err = fillData(data, "photo", dal.GetQueryByCtx(ctx).Photo.WithContext(ctx).Find, err)
	err = fillData(data, "post", dal.GetQueryByCtx(ctx).Post.WithContext(ctx).Find, err)
	err = fillData(data, "post_category", dal.GetQueryByCtx(ctx).PostCategory.WithContext(ctx).Find, err)
	err = fillData(data, "post_tag", dal.GetQueryByCtx(ctx).PostTag.WithContext(ctx).Find, err)
	err = fillData(data, "theme_setting", dal.GetQueryByCtx(ctx).ThemeSetting.WithContext(ctx).Find, err)
	err = fillData(data, "user", dal.GetQueryByCtx(ctx).User.WithContext(ctx).Find, err)
	if err != nil {
		return nil, xerr.NoType.Wrap(err).WithStatus(xerr.StatusInternalServerError)
	}

	backupFilename := consts.SonicDataExportPrefix + time.Now().Format("2006-01-02-15-04-05") + util.GenUUIDWithOutDash() + ".json"

	backupFilePath := config.DataExportDir
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		err = os.MkdirAll(backupFilePath, os.ModePerm)
		if err != nil {
			return nil, xerr.NoType.Wrap(err).WithMsg("create dir err")
		}
	} else if err != nil {
		return nil, xerr.NoType.Wrap(err).WithMsg("get fileInfo")
	}
	backupFile := filepath.Join(backupFilePath, backupFilename)

	file, err := os.Create(backupFile)
	if err != nil {
		return nil, xerr.NoType.Wrap(err).WithMsg("get fileInfo")
	}
	defer file.Close()
	content, err := json.Marshal(data)
	if err != nil {
		return nil, xerr.NoType.Wrap(err).WithStatus(xerr.StatusInternalServerError).WithMsg("json marshal err")
	}
	_, err = file.Write(content)
	if err != nil {
		return nil, xerr.NoType.Wrap(err).WithMsg("write to file err")
	}
	return b.buildBackupDTO(ctx, string(service.JsonData), filepath.Join(backupFilePath, backupFilename))
}

func (b *backupServiceImpl) ExportMarkdown(ctx context.Context, needFrontMatter bool) (*dto.BackupDTO, error) {
	// TODO
	return nil, nil
}

func (b *backupServiceImpl) buildBackupDTO(ctx context.Context, baseBackupURL string, backupFilePath string) (*dto.BackupDTO, error) {
	backupDTO := &dto.BackupDTO{}
	backupFilename := filepath.Base(backupFilePath)
	downloadLink, err := b.buildDownloadURL(ctx, baseBackupURL, backupFilePath)
	if err != nil {
		return nil, err
	}
	backupDTO.DownloadLink = downloadLink
	backupDTO.Filename = backupFilename
	fileInfo, err := os.Stat(backupFilePath)
	if err != nil {
		return nil, xerr.NoType.Wrap(err).WithMsg("Failed to access file")
	}
	backupDTO.UpdateTime = fileInfo.ModTime().UnixMilli()
	backupDTO.FileSize = fileInfo.Size()
	return backupDTO, nil
}

func (b *backupServiceImpl) buildDownloadURL(ctx context.Context, baseBackupURL, backupFilePath string) (string, error) {
	backupFileURL := baseBackupURL + "/" + filepath.Base(backupFilePath)

	oneTimeToken := b.OneTimeTokenService.Create(backupFileURL)
	blogURL, err := b.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return "", err
	}
	return util.CompositeURL(blogURL, backupFileURL+"?"+consts.OneTimeTokenQueryName+"="+oneTimeToken), nil
}

// TODO when refactor dal ,refactor this method
func fillData(dataMap map[string]interface{}, item string, f interface{}, preErr error) error {
	if preErr != nil {
		return preErr
	}
	fv := reflect.ValueOf(f)
	params := []reflect.Value{}
	rs := fv.Call(params)
	data := rs[0]
	err := rs[1]
	if !err.IsNil() {
		return WrapDBErr(err.Interface().(error))
	}
	if !data.IsNil() {
		dataMap[item] = data.Interface()
	}
	return nil
}

func (b *backupServiceImpl) ListToBackupItems(ctx context.Context) ([]string, error) {
	dirEntrys, err := os.ReadDir(b.Config.Sonic.WorkDir)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("read work dir err")
	}

	result := make([]string, 0)
	for _, dirEntry := range dirEntrys {
		result = append(result, dirEntry.Name())
	}
	return result, nil
}
