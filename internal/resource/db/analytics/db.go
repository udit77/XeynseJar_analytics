package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xeynse/XeynseJar_analytics/internal/config"
)

type Resource interface {
	GetClient() *sqlx.DB
	PreparexFatal(query string) *sqlx.Stmt
	PrepareNamedFatal(query string) *sqlx.NamedStmt
}

type resource struct {
	db *sqlx.DB
}

func New(config *config.Config) (Resource, error) {
	dBConfig := config.DataBase.XeynsJar.JarStatus
	dataSource := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", dBConfig.UserName, dBConfig.Password, dBConfig.Host, dBConfig.Name)
	dbConn, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return &resource{}, err
	}
	return &resource{
		db: dbConn,
	}, err
}

func (r resource) GetClient() *sqlx.DB {
	return r.db
}

func (r resource) PreparexFatal(query string) *sqlx.Stmt {
	stmt, err := r.db.Preparex(query)
	if err != nil {
		log.Fatal("[PreparexFatal] Fatal occurred for query", query)
	}
	return stmt
}

func (r resource) PrepareNamedFatal(query string) *sqlx.NamedStmt {
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		log.Fatal("[PrepareNamedFatal] Fatal occurred for query", query, err)
	}
	return stmt
}
