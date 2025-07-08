package db

import (
	"fmt"
	logger_pkg "github.com/IntelXLabs-LLC/go-common/logger"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/IntelXLabs-LLC/go-common/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

// DatabaseConnection represents a connection to a database.
type DatabaseConnection struct {
	db *gorm.DB
}

// Disconnect closes the underlying database connection associated with the DatabaseConnection instance.
// It first obtains the underlying raw `*sql.DB` object by calling `d.db.DB()`, and then calls `Close()` on it.
// If any error occurs during obtaining the raw database connection or while closing it, the error is returned.
func (d *DatabaseConnection) Disconnect() error {
	sqlDb, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

// Raw returns the raw Gorm DB object associated with the DatabaseConnection instance.
func (d *DatabaseConnection) Raw() *gorm.DB {
	return d.db
}

// ConnectDB establishes a connection to the database using the provided configuration.
//
// The function takes a `config.DatabaseConfig` object as input, which contains the necessary information to establish the connection.
// It constructs a data source name (DSN) string based on the configuration properties.
//
// The function then opens the connection to the database using the constructed DSN and the Gorm package.
// It returns a `*DatabaseConnection` object, which encapsulates the Gorm DB object.
//
// If an error occurs during the connection process, it will be returned along with the `*DatabaseConnection` object as `nil`.
// Any subsequent Gorm operations using the returned `*DatabaseConnection` object should handle potential errors appropriately.
//
// Example usage:
//
//	config := config.DatabaseConfig{
//	    Host:     "localhost",
//	    Port:     "5432",
//	    User:     "postgres",
//	    Password: "password",
//	    DBName:   "mydb",
//	    SSLMode:  "disable",
//	}
//	conn, err := ConnectDB(config)
//	if err != nil {
//	    fmt.Println("Failed to connect to the database:", err)
//	} else {
//	    defer conn.Close()
//	    // Perform database operations using `conn`
//	}
func ConnectDB(config config.DatabaseConfig) (*DatabaseConnection, error) {
	// DSN represents data source name for the database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	// Open the connection to the database
	// db, err := sql.Open("postgres", dsn) // lib/pq version
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy:         nil,
		FullSaveAssociations:   false,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		}),
		PrepareStmt:          true,
		DisableAutomaticPing: false,
		CreateBatchSize:      10,
	})
	if err != nil {
		panic(err)
	}

	// Create sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Print a message to the console if the connection is successful
	logger_pkg.Info("db connected successfully")

	// Return DatabaseConnection
	return &DatabaseConnection{
		db: db,
	}, err
}
