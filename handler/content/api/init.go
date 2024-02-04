package api

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewArchiveHandler,
		NewCategoryHandler,
		NewJournalHandler,
		NewLinkHandler,
		NewPostHandler,
		NewSheetHandler,
		NewOptionHandler,
		NewPhotoHandler,
		NewCommentHandler,
	)
}
