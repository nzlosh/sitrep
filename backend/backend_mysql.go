package backend_mysql

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "log"
)

func Name() string {
    return "MySQL"
}

type Impl struct {
    DB *gorm.DB
}

func (i *Impl) InitDB(cxn string) {
    var err error
    i.DB, err = gorm.Open("mysql", cxn)
    if err != nil {
        log.Fatalf("Got error when connecting to database, the error is '%v'", err)
    }
    i.DB.LogMode(true)
}
