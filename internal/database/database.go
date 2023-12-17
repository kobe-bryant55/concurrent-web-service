package database

import (
	"database/sql"
)

type SQLStore interface {
	InitDB()
	GetInstance() *sql.DB
}

type StoreOpts struct {
	user     string
	name     string
	password string
	port     string
}

type StoreOptsFunc func(*StoreOpts)

func WithUser(u string) StoreOptsFunc {
	return func(opts *StoreOpts) {
		opts.user = u
	}
}

func WithPassword(p string) StoreOptsFunc {
	return func(opts *StoreOpts) {
		opts.password = p
	}
}

func WithName(n string) StoreOptsFunc {
	return func(opts *StoreOpts) {
		opts.name = n
	}
}

func WithPort(p string) StoreOptsFunc {
	return func(opts *StoreOpts) {
		opts.port = p
	}
}
