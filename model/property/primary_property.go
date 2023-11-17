package property

import (
	"reflect"

	"github.com/go-sonic/sonic/consts"
)

var (
	IsInstalled = Property{
		KeyValue:     "is_installed",
		DefaultValue: false,
		Kind:         reflect.Bool,
	}
	Theme = Property{
		KeyValue:     "theme",
		DefaultValue: consts.DefaultThemeID,
		Kind:         reflect.String,
	}
	BirthDay = Property{
		KeyValue:     "birthday",
		DefaultValue: int64(0),
		Kind:         reflect.Int64,
	}
	DefaultMenuTeam = Property{
		KeyValue:     "default_menu_team",
		DefaultValue: "",
		Kind:         reflect.String,
	}
)
