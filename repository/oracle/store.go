package oracle

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"test/conf"
	"test/logger"
	"test/repository"

	_ "github.com/godror/godror"
	"github.com/spf13/viper"
)

var store repository.Store
var storeOnce sync.Once

type Store struct {
	DB     *sql.DB
	Logger logger.Logger
}

// SharedStore return global or single instance of firebase connection (bounded in sync once)
func SharedStore(logger logger.Logger) repository.Store {
	storeOnce.Do(func() {
		var err error
		store, err = NewStore(logger)
		if err != nil {
			panic(err)
		}
	})

	return store
}

// NewStore creates a new instance of the Store
func NewStore(logger logger.Logger) (*Store, error) {
	// Replace these values with your actual Oracle database connection details
	dbUser := viper.Get(conf.DbUser).(string)
	dbPassword := viper.Get(conf.DbPassword).(string)
	dbHost := viper.Get(conf.DbHost).(string)
	dbPort := viper.Get(conf.DbPort).(string)
	dbName := viper.Get(conf.DbName).(string)

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s/%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open a database connection
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is valid
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the Oracle database")

	return &Store{
		Logger: logger,
		DB:     db,
	}, nil
}
