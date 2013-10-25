package cache

import (
	"github.com/efimbakulin/email-distribution/dao"
	"text/template"
)

type Templates struct {
	templatesDao *dao.Templates
	cache        map[string]*template.Template
}

var (
	tplCache *Templates
)

func (self *Templates) Get(lang string) (*template.Template, error) {
	if self.cache[lang] != nil {
		return self.cache[lang], nil
	}
	tpl, err := self.templatesDao.LoadByLang(lang)
	if err != nil {
		return nil, nil
	}
	self.cache[lang] = template.Must(template.New(string(lang)).Parse(tpl))

	return self.cache[lang], nil
}

func NewTemplateCache(connectionString string) *Templates {
	if nil == tplCache {
		tplCache = &Templates{
			templatesDao: dao.NewTemplatesDao(connectionString),
			cache:        make(map[string]*template.Template),
		}

	}
	return tplCache
}
