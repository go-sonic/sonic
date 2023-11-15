package impl

import (
	"github.com/go-sonic/sonic/injection"
	"github.com/go-sonic/sonic/service/storage"
)

func init() {
	injection.Provide(
		NewAdminService,
		NewAttachmentService,
		NewAuthenticateService,
		NewBackUpService,
		NewBaseCommentService,
		NewBasePostService,
		NewCategoryService,
		NewEmailService,
		NewInstallService,
		NewJournalService,
		NewLinkService,
		NewJournalCommentService,
		NewLogService,
		NewMenuService,
		NewMetaService,
		NewBaseMFAService,
		NewTwoFactorTOTPMFAService,
		NewOneTimeTokenService,
		NewOptionService,
		NewClientOptionService,
		NewPhotoService,
		NewPostService,
		NewPostCategoryService,
		NewPostCommentService,
		NewPostTagService,
		NewSheetService,
		NewSheetCommentService,
		NewStatisticService,
		NewTagService,
		NewThemeService,
		NewUserService,
		NewExportImport,
		storage.NewFileStorageComposite,
		NewApplicationPasswordService,
	)
}
