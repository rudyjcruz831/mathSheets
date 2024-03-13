package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dataSources struct {
	DB *gorm.DB
	//RedisClient *redis.Client
	//StorageClient *storage.Client

	sqlDB *sql.DB
}

// Itilizing data source connections to field in data source
func initDS() (*dataSources, error) {
	fmt.Println("Initilizing data source")

	// load env variables for postgres
	dsn := fmt.Sprintf("host=%s dbname=%s port=%s user=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"))

	log.Printf("Connecting to Postgresql")
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		fmt.Print("Unable to connect to db \n")
		panic(err)
	}
	// this returns gorms own interface to use to ping in next lines of code still dont understand
	// this website can help udnerstand more if need = https://gorm.io/docs/generic_interface.html
	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get sqlDB from db.DB() error : %v", err)
		return nil, err
	}

	// verify database connection is working by ping database
	if err := sqlDB.Ping(); err != nil {
		for i := 0; i < 3; i++ {
			fmt.Println("Ping")
		}
		return nil, err
	}

	// fmt.Print("database successfully configured\n")
	return &dataSources{DB: db, sqlDB: sqlDB}, nil

}

func (d *dataSources) Close() error {
	if err := d.sqlDB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	return nil
}
