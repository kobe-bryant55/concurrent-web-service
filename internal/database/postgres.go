package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	singleton *postgresStore
	storeOnce sync.Once
	initOnce  sync.Once
)

type postgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(opts ...StoreOptsFunc) SQLStore {
	storeOnce.Do(func() {
		o := postgresStoreDefaultOpts()
		for _, fn := range opts {
			fn(&o)
		}

		conStr := fmt.Sprintf(`user=%s dbname=%s password=%s port=%s sslmode=disable`, o.user, o.name, o.password, o.port)

		db, err := sql.Open("postgres", conStr)
		if err != nil {
			log.Fatal(err.Error())
		}

		singleton = &postgresStore{DB: db}

		if err := db.Ping(); err != nil {
			log.Fatal("Connection Error: " + err.Error())
		}
	})

	return singleton
}

func (s postgresStore) GetInstance() *sql.DB {
	return s.DB
}

func (s postgresStore) InitDB() {
	initOnce.Do(func() {
		s.createStatusEnums()
		s.createTasksTable()
	})
}

func (s postgresStore) createStatusEnums() {
	query := `DO $$ BEGIN
		IF to_regtype('status') IS NULL THEN
		CREATE TYPE status AS ENUM('active', 'passive');
		END IF;
		END $$;
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal("tasks status enums creation failed: ", err.Error())
	}
}

func (s postgresStore) createTasksTable() {
	query := `CREATE TABLE IF NOT EXISTS tasks (
    	id				serial PRIMARY KEY,
    	description 	varchar(500) NOT NULL, 
    	title 		   	varchar(21) NOT NULL, 
		status_state    status
	)
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal("tasks table creation failed: ", err.Error())
	}
}

func postgresStoreDefaultOpts() StoreOpts {
	return StoreOpts{
		user:     "development",
		name:     "development",
		password: "development",
		port:     "5432",
	}
}
