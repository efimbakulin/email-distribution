package main

import (
	"github.com/efimbakulin/email-distribution/dao"
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
	cache    = make(map[CacheKey]string)
)

func (self *Cache) Get(key CacheKey) (string, error) {
	if cache[key] != "" {
		log.Printf("Found in cache")
		return cache[key], nil
	}
	template, err := self.templatesDao.LoadTemplate(key.TemplateId)
	if err != nil {
		return "", nil
	}
	letter, err := self.lettersDao.LoadLetter(key.LetterId)
	if err != nil {
		return "", nil
	}
	_ = template
	cache[key] = letter.Body

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
