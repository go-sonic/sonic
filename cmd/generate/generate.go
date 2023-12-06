package main

import (
	"go.uber.org/fx"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/log"
)

// generate code
func main() {
	var DB *gorm.DB
	_ = fx.New(
		fx.Provide(log.NewLogger),
		fx.Provide(dal.NewGormDB),
		fx.Provide(log.NewGormLogger),
		fx.Provide(config.NewConfig),
		fx.Populate(&DB),
	)
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		Mode:         gen.WithDefaultQuery,
		OutPath:      "./dal",
		ModelPkgPath: "./model/entity",
		/* Mode: gen.WithoutContext,*/
		// if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary, or it will panic
	g.UseDB(DB)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	g.ApplyBasic(g.GenerateModel("attachment", gen.FieldType("type", "consts.AttachmentType")),
		g.GenerateModel("category", gen.FieldType("type", "consts.CategoryType")),
		g.GenerateModel("comment", gen.FieldType("type", "consts.CommentType"), gen.FieldType("status", "consts.CommentStatus")),
		g.GenerateModel("comment_black"),
		g.GenerateModel("journal", gen.FieldType("type", "consts.JournalType")),
		g.GenerateModel("link"),
		g.GenerateModel("log", gen.FieldType("type", "consts.LogType")),
		g.GenerateModel("menu"),
		g.GenerateModelAs("meta", "Meta", gen.FieldType("type", "consts.MetaType")),
		g.GenerateModel("option", gen.FieldType("type", "consts.OptionType")),
		g.GenerateModel("photo"),
		g.GenerateModel("post", gen.FieldType("type", "consts.PostType"), gen.FieldType("status", "consts.PostStatus"), gen.FieldType("editor_type", "consts.EditorType")),
		g.GenerateModel("post_category"),
		g.GenerateModel("post_tag"),
		g.GenerateModel("tag"),
		g.GenerateModel("theme_setting"),
		g.GenerateModel("user", gen.FieldType("mfa_type", "consts.MFAType")),
		//g.GenerateModel("application_password"),
		g.GenerateModel("scrap_page"),
	)

	// apply diy interfaces on structs or table models
	// g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

	// execute the action of code generation
	g.Execute()
}
