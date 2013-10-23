package cache

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
	stringBuilder, _ = connstring.CreateBuilder(connstring.ConnectionStringPg)
)

func init() {
	stringBuilder.Address(Addr)
	stringBuilder.Port(Port)
	stringBuilder.Username(Username)
	stringBuilder.Password(Password)
	stringBuilder.Dbname(Dbname)
}

func Test_Templates(test *testing.T) {
	templates := NewTemplateCache(stringBuilder.Build())
	template, err := templates.Get(1)
	if err != nil {
		test.Fatal(err)
	}
	test.Log(template)
}

func Test_Letters(test *testing.T) {
	letters := NewLetterCache(stringBuilder.Build())
	letter, err := letters.Get(1)
	if err != nil {
		test.Fatal(err)
	}
	test.Log(letter.Subj)
	test.Log(letter.Body)
}
