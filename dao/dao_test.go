package dao

import (
	"github.com/efimbakulin/connection-string-builder"
	"testing"
)

const (
	Addr     = "unstable-ebakulin.2reallife.net"
	Port     = 56432
	Username = "hv"
	Password = "oops"
	Dbname   = "sa_data"
)

var (
	testEmails       = []string{"1@mail.ru", "2@mail.ru"}
	stringBuilder, _ = builder.New(builder.ConnectionStringPg)
)

func init() {
	stringBuilder.Address(Addr)
	stringBuilder.Port(Port)
	stringBuilder.Username(Username)
	stringBuilder.Password(Password)
	stringBuilder.Dbname(Dbname)
}

func TestEmailMark(test *testing.T) {
	test.Log(stringBuilder.Build())
	emailsDao := NewEmailsDao(stringBuilder.Build())
	_, err := emailsDao.MarkInvalid(testEmails)
	if err != nil {
		test.Fatal(err)
	}
}

func TestTemplateLoad(test *testing.T) {
	dao := NewTemplatesDao(stringBuilder.Build())
	body, err := dao.LoadTemplate(1)
	if err != nil {
		test.Fatal(err)
	}

	test.Log(body)
}

func TestLetterLoad(test *testing.T) {
	dao := NewLettersDao(stringBuilder.Build())
	body, err := dao.LoadLetter(1)
	if err != nil {
		test.Fatal(err)
	}

	test.Log(body.Body)
	test.Log(body.Subj)
}
