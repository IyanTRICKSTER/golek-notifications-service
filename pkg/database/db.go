package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	Host     string
	Username string
	Password string
	DbName   string
	DbPort   string
	conn     *gorm.DB
}

func (d *Database) Connect() error {
	if d.conn == nil {
		dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", d.Username, d.Password, d.Host, d.DbPort, d.DbName)
		conn, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
		if err != nil {
			log.Println("Can't establish a connection: ", err)
			return err
		}
		log.Println("Connection has been established")
		d.conn = conn
	} else {
		log.Println("Connection already established")
	}

	return nil
}

func (d *Database) GetConnection() *gorm.DB {
	return d.conn
}

func (d *Database) MigrateTable(entity interface{}) error {
	err := d.conn.AutoMigrate(&entity)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (d *Database) DropTable(entity interface{}) error {

	if !d.conn.Migrator().HasTable(entity) {
		return errors.New("table is not exists")
	}

	err := d.conn.Migrator().DropTable(entity)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

//func (d *Database) Seed(entities []entities.User) []error {
//	var errs []error
//	for _, entity := range entities {
//		err := d.conn.Create(&entity).Error
//		errs = append(errs, err)
//	}
//	return errs
//}
