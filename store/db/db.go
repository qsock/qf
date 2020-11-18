package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var x map[string]*sql.DB = make(map[string]*sql.DB)

var ErrNoDb = errors.New("no db")

type Config struct {
	Addr            string `toml:"addr" json:"addr"`
	User            string `toml:"user" json:"user"`
	Pwd             string `toml:"pwd" json:"pwd"`
	Db              string `toml:"db" json:"db"`
	Options         string `toml:"options" json:"options"`
	MaxOpenConns    int    `toml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns    int    `toml:"max_idle_conns" json:"max_idle_conns"`
	MaxConnLifeTime int    `toml:"max_conn_life_time" json:"max_conn_life_time"`
}

func (c Config) Dsn() string {
	if c.Options == "" {
		format := "%s:%s@tcp(%s)/%s?charset=utf8mb4&timeout=5s"
		return fmt.Sprintf(format, c.User, c.Pwd, c.Addr, c.Db)
	} else {
		format := "%s:%s@tcp(%s)/%s?charset=utf8mb4&timeout=5s&%s"
		return fmt.Sprintf(format, c.User, c.Pwd, c.Addr, c.Db, c.Options)
	}
}

func Add(name string, c Config) error {
	db, err := sql.Open("mysql", c.Dsn())
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)

	if c.MaxConnLifeTime <= 0 {
		c.MaxConnLifeTime = 600 // Default 10 minutes
	}
	db.SetConnMaxLifetime(time.Duration(c.MaxConnLifeTime) * time.Second)

	_, ok := x[name]
	if ok {
		return errors.New("db exists")
	}
	x[name] = db
	return nil
}

func GetDB(name string) (*sql.DB, error) {
	db, ok := x[name]
	if ok {
		return db, nil
	}
	return nil, ErrNoDb
}
