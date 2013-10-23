package dao

import (
	"fmt"
	"strings"
)

const (
	QueryMarkEmailsInvalid = "select delivery.mark_emails_unused('{%s}'::varchar[])"
)

type Emails struct {
	*BaseDao
}

func (self *Emails) MarkInvalid(emails []string) (int64, error) {
	return self.executeQuery(fmt.Sprintf(QueryMarkEmailsInvalid, strings.Join(emails, ", ")))
}

func NewEmailsDao(connectionString string) *Emails {
	return &Emails{&BaseDao{connectionString}}
}
