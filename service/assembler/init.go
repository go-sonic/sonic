package assembler

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewBasePostAssembler,
		NewPostAssembler,
		NewSheetAssembler,
		NewBaseCommentAssembler,
		NewPostCommentAssembler,
		NewJournalCommentAssembler,
		NewSheetCommentAssembler,
	)
}
