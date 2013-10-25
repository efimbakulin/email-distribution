package dao

import (
	"fmt"
	"strings"
)

const (
	QueryMarkEmailsInvalid = "select delivery.mark_emails_unused('{%s}'::varchar[])"
	QueryGetIdByEmail      = "select ref_id || '_' || ref_nick from auth.users where email = '%s'"
)

type Emails struct {
	*BaseDao
}

func (self *Emails) MarkInvalid(emails []string) (int64, error) {
	return self.executeQuery(fmt.Sprintf(QueryMarkEmailsInvalid, strings.Join(emails, ", ")))
}

func (self *Emails) GetId(email string) (int64, error) {
	return 1, nil
}

func NewEmailsDao(connectionString string) *Emails {
	return &Emails{&BaseDao{connectionString}}
}
