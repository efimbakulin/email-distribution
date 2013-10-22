package main

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
	stringBuilder, _ = connstring.CreateBuilder(connstring.ConnectionStringPg)
)

func init() {
	stringBuilder.Address(Addr)
	stringBuilder.Port(Port)
	stringBuilder.Username(Username)
	stringBuilder.Password(Password)
	stringBuilder.Dbname(Dbname)
}

func TestCache(test *testing.T) {
	theCache := NewCache(stringBuilder.Build())
	template, err := theCache.Get(CacheKey{1, 2})
	if err != nil {
		test.Fatal(err)
	}
	test.Log(template)
	template, err = theCache.Get(CacheKey{1, 1})
	template, err = theCache.Get(CacheKey{1, 2})
}
