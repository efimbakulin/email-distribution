package cache

import (
	"github.com/efimbakulin/email-distribution/dao"
)

type LetterInfo struct {
	Subj string
	Body string
}

type Letters struct {
	dao   *dao.Letters
	cache map[int64]*LetterInfo
}

var (
	letterCache *Letters
)

func (self *Letters) Get(key int64) (*LetterInfo, error) {
	if self.cache[key] != nil {
		return self.cache[key], nil
	}
	var err error
	letter := &LetterInfo{}
	letter.Subj, letter.Body, err = self.dao.LoadLetter(key)
	if err != nil {
		return nil, nil
	}
	self.cache[key] = letter

	return letter, nil
}

func NewLetterCache(connectionString string) *Letters {
	if nil == letterCache {
		letterCache = &Letters{
			dao:   dao.NewLettersDao(connectionString),
			cache: make(map[int64]*LetterInfo),
		}

	}
	return letterCache
}
