package admin

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewAdminHandler,
		NewAttachmentHandler,
		NewCategoryHandler,
		NewBackupHandler,
		NewInstallHandler,
		NewJournalHandler,
		NewJournalCommentHandler,
		NewLinkHandler,
		NewLogHandler,
		NewMenuHandler,
		NewOptionHandler,
		NewPhotoHandler,
		NewPostHandler,
		NewPostCommentHandler,
		NewSheetHandler,
		NewSheetCommentHandler,
		NewStatisticHandler,
		NewTagHandler,
		NewThemeHandler,
		NewUserHandler,
		NewEmailHandler,
		NewApplicationPasswordHandler,
	)
}
