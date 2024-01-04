package content

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewIndexHandler,
		NewFeedHandler,
		NewArchiveHandler,
		NewViewHandler,
		NewCategoryHandler,
		NewSheetHandler,
		NewTagHandler,
		NewLinkHandler,
		NewPhotoHandler,
		NewJournalHandler,
		NewSearchHandler,
		NewScrapHandler,
	)
}
