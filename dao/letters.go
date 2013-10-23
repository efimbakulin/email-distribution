package dao

import (
	"fmt"
	"github.com/lxn/go-pgsql"
)

const (
	QueryLoadLetter = "select subject, body from delivery.deliveries where delivery_id = @id;"
)

type Letters struct {
	*BaseDao
}

func (self *Letters) LoadLetter(letter_id int64) (Subject string, Body string, err error) {
	connection, err := self.getConnection()
	if err != nil {
		return Subject, Body, err
	}

	parameter := pgsql.NewParameter("@id", pgsql.Bigint)
	err = parameter.SetValue(letter_id)
	recordSet, err := connection.Query(QueryLoadLetter, parameter)
	if err != nil {
		return Subject, Body, err
	}
	defer recordSet.Close()
	fetched, err := recordSet.FetchNext()
	if err != nil {
		return Subject, Body, err
	}
	if !fetched {
		return Subject, Body, fmt.Errorf("Letter not found")
	}

	err = recordSet.Scan(&Subject, &Body)
	if err != nil {
		return Subject, Body, err
	}

	return Subject, Body, nil
}

func NewLettersDao(connectionString string) *Letters {
	return &Letters{&BaseDao{connectionString}}
}
