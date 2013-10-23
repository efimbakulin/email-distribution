package cache

import (
	"github.com/efimbakulin/email-distribution/dao"
	"text/template"
)

type Templates struct {
	templatesDao *dao.Templates
	cache        map[int64]*template.Template
}

var (
	tplCache *Templates
)

func (self *Templates) Get(key int64) (*template.Template, error) {
	if self.cache[key] != nil {
		return self.cache[key], nil
	}
	tpl, err := self.templatesDao.LoadTemplate(key)
	if err != nil {
		return nil, nil
	}
	self.cache[key] = template.Must(template.New(string(key)).Parse(tpl))

	return self.cache[key], nil
}

func NewTemplateCache(connectionString string) *Templates {
	if nil == tplCache {
		tplCache = &Templates{
			templatesDao: dao.NewTemplatesDao(connectionString),
			cache:        make(map[int64]*template.Template),
		}

	}
	return tplCache
}
