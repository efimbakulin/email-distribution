package dao

import (
	"fmt"
	"strings"
)

const (
	QueryMarkEmailsInvalid = "select delivery.mark_emails_unused('{%s}'::varchar[])"
)

type EmailsDao struct {
	*BaseDao
}

func (self *EmailsDao) MarkInvalid(emails []string) (int64, error) {
	return self.executeQuery(fmt.Sprintf(QueryMarkEmailsInvalid, strings.Join(emails, ", ")))
}

func NewEmailsDao(connectionString string) *EmailsDao {
	return &EmailsDao{&BaseDao{connectionString}}
}
