package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/phper-go/frame/func/conv"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/phper-go/frame/func/conv"
)

type PDO interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type DB struct {
	db           *sql.DB
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	Dbname       string
	Charset      string
	Timeout      int
	MaxLifetime  int
	MaxIdleConns int
	MaxOpenConns int
}

func (this *DB) SqlDB() *sql.DB {
	return this.db
}

func (this *DB) Init() error {

	if this.Host == "" {
		return errors.New("host empty")
	}
	if this.Port == 0 {
		return errors.New("port empty")
	}
	if this.User == "" {
		return errors.New("user empty")
	}
	if this.Password == "" {
		return errors.New("password empty")
	}
	if this.Dbname == "" {
		return errors.New("dbname empty")
	}
	if this.Host == "" {
		return errors.New("host empty")
	}
	if this.Driver == "" {
		this.Driver = "mysql"
	}
	conn, err := sql.Open(this.Driver, this.ConnString())
	if err == nil {

		this.db = conn
		if this.MaxLifetime == 0 {
			this.MaxLifetime = 3600
		}
		this.db.SetConnMaxLifetime(time.Duration(this.MaxLifetime) * time.Second)
		if this.MaxIdleConns == 0 {
			this.MaxIdleConns = 1000
		}
		this.db.SetMaxIdleConns(this.MaxIdleConns)
		if this.MaxOpenConns == 0 {
			this.MaxOpenConns = 2000
		}
		this.db.SetMaxOpenConns(this.MaxOpenConns)
		return this.db.Ping()
	}
	return err
}

func (this *DB) Close() {
	this.db.Close()
}

func (this *DB) ConnString() string {
	//user:password@tcp(ip:port)/database
	if this.Charset == "" {
		this.Charset = "utf8"
	}
	if this.Timeout == 0 {
		this.Timeout = 30
	}
	connString := this.User + ":" + this.Password + "@tcp(" + this.Host + ":" + conv.String(this.Port) + ")/" + this.Dbname
	connString += "?charset=" + this.Charset + "&timeout=" + conv.String(this.Timeout) + "s"
	return connString
}

func (this *DB) Cmd(ctx ...context.Context) *DBCommand {

	return Cmd(this.db, ctx...)
}

func Cmd(pdo PDO, ctx ...context.Context) *DBCommand {

	command := &DBCommand{}
	command.SetPDO(pdo)
	if len(ctx) > 0 {
		command.SetContext(ctx[0])
	}
	return command
}

func Commit(db *sql.DB, txRun func(tx *sql.Tx, ctx ...context.Context) error) error {

	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		if err := txRun(tx * sql.Tx); err != nil {
			return tx.Rollback()
		} else {
			return tx.Commit()
		}
	}

}
