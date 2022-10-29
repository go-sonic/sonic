package extension

import (
	"math/rand"
	"time"

	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type toolExtension struct {
	Template *template.Template
}

func RegisterToolFunc(template *template.Template) {
	t := &toolExtension{
		Template: template,
	}
	t.addRainbow()
	t.addRandom()
}

// addRainbow 彩虹分页算法
func (t *toolExtension) addRainbow() {
	rainbowPage := func(page, total, display int) []int {
		return util.RainbowPage(page, total, display)
	}
	t.Template.AddFunc("rainbowPage", rainbowPage)
}

func (t *toolExtension) addRandom() {
	rand.Seed(time.Now().UnixNano())
	random := func(min, max int) int {
		return min + rand.Intn(max-min)
	}
	t.Template.AddFunc("random", random)
}
