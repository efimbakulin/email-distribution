package dao

import (
	"github.com/lxn/go-pgsql"
)

type BaseDao struct {
	connectionString string
}

func (self *BaseDao) getConnection() (*pgsql.Conn, error) {
	connection, err := pgsql.Connect(self.connectionString, pgsql.LogNothing)
	if nil != err {
		return nil, err
	}
	return connection, err
}

func (self *BaseDao) executeQuery(query string) (int64, error) {
	connection, err := self.getConnection()
	if err != nil {
		return 0, err
	}
	defer connection.Close()
	return connection.Execute(query)
}

func NewBaseDao(connectionString string) (*BaseDao, error) {
	var err error
	instance := new(BaseDao)
	instance.connectionString = connectionString
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (self *BaseDao) Close() error {
	return nil
}
