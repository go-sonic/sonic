package dto

import "github.com/go-sonic/sonic/consts"

type ThemeProperty struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name"`
	Website        string                     `json:"website"`
	Branch         string                     `json:"branch"`
	Repo           string                     `json:"repo"`
	UpdateStrategy consts.ThemeUpdateStrategy `json:"updateStrategy"`
	Description    string                     `json:"description"`
	Logo           string                     `json:"logo"`
	Version        string                     `json:"version"`
	Require        string                     `json:"require"`
	Author         ThemeAuthor                `json:"author"`
	ThemePath      string                     `json:"themePath"`
	FolderName     string                     `json:"folderName"`
	HasOptions     bool                       `json:"hasOptions"`
	Activated      bool                       `json:"activated"`
	ScreenShots    string                     `json:"screenshots"`
	PostMetaField  []string                   `json:"postMetaField"`
	SheetMetaField []string                   `json:"sheetMetaField"`
}

type ThemeAuthor struct {
	Name    string `json:"name"`
	Website string `json:"website"`
	Avatar  string `json:"avatar"`
}

type ThemeFile struct {
	Name     string       `json:"name"`
	Path     string       `json:"path"`
	IsFile   bool         `json:"isFile"`
	Editable bool         `json:"editable"`
	Node     []*ThemeFile `json:"node"`
}

type ThemeConfigGroup struct {
	Name    string                      `json:"name"`
	Label   string                      `json:"label"`
	Items   []*ThemeConfigItem          `json:"items" yaml:"-"`
	ItemMap map[string]*ThemeConfigItem `json:"-" yaml:"items"`
}

type ThemeConfigItem struct {
	Name         string                      `json:"name"`
	Label        string                      `json:"label"`
	InputType    consts.ThemeConfigInputType `json:"type" yaml:"type"`
	DataType     consts.ThemeConfigDataType  `json:"dataType" yaml:"data-type"`
	DefaultValue interface{}                 `json:"defaultValue" yaml:"default"`
	PlaceHolder  string                      `json:"placeHolder"`
	Description  string                      `json:"description"`
	Options      []*ThemeConfigOption        `json:"options"`
}

type ThemeConfigOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}
