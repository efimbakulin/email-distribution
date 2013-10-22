package main

import (
	"github.com/efimbakulin/email-distribution/dao"
	"html/template"
	"log"
)

type Cache struct {
	templatesDao *dao.TemplatesDao
	lettersDao   *dao.LettersDao
}

type CacheKey struct {
	LetterId   int64
	TemplateId int64
}

var (
	instance *Cache
	cache    = make(map[CacheKey]*template.Template)
)

func (self *Cache) Get(key CacheKey) (*template.Template, error) {
	if cache[key] != nil {
		log.Printf("Found in cache")
		return cache[key], nil
	}
	tpl, err := self.templatesDao.LoadTemplate(key.TemplateId)
	if err != nil {
		return nil, nil
	}
	letter, err := dao.LetterInfo{"Letter body", "Subject"}, nil
	if err != nil {
		return nil, nil
	}
	_ = letter
	cache[key] = template.Must(template.New(string(key.TemplateId)).Parse(tpl))

	return cache[key], nil
}

func NewCache(connectionString string) *Cache {
	if nil == instance {
		instance = &Cache{
			lettersDao:   dao.NewLettersDao(connectionString),
			templatesDao: dao.NewTemplatesDao(connectionString),
		}

	}
	return instance
}
