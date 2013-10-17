package dao

import (
	"fmt"
	"github.com/lxn/go-pgsql"
)

const (
	QueryLoadLetter = "select subject, body from delivery.deliveries where delivery_id = @id;"
)

type LettersDao struct {
	*BaseDao
}

type LetterInfo struct {
	Body string
	Subj string
}

func (self *LettersDao) LoadLetter(letter_id int64) (*LetterInfo, error) {
	connection, err := self.getConnection()
	if err != nil {
		return nil, err
	}

	parameter := pgsql.NewParameter("@id", pgsql.Bigint)
	err = parameter.SetValue(letter_id)
	recordSet, err := connection.Query(QueryLoadLetter, parameter)
	if err != nil {
		return nil, err
	}
	defer recordSet.Close()
	fetched, err := recordSet.FetchNext()
	if err != nil {
		return nil, err
	}
	if !fetched {
		return nil, fmt.Errorf("Letter not found")
	}

	result := &LetterInfo{}
	err = recordSet.Scan(&result.Subj, &result.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewLettersDao(connectionString string) *LettersDao {
	return &LettersDao{&BaseDao{connectionString}}
}
