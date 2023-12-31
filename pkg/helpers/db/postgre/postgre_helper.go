package src

import (
	"database/sql"
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"time"
)

func OpenPostgresConnection(dbConfig *structures.DbConfig) (*sql.DB, error) {
	// Open up database connection.
	db, err := sql.Open("postgres", getConnectionString(dbConfig))

	if err != nil {
		log.Fatal(err)
	}

	maxPoolConnections, err := strconv.Atoi(dbConfig.MaxPoolConnections)

	if err == nil {
		db.SetMaxOpenConns(maxPoolConnections)
	}

	maxIdlePoolConnections, err := strconv.Atoi(dbConfig.MaxIdlePoolConnections)

	if err == nil {
		db.SetMaxIdleConns(maxIdlePoolConnections)
	}

	connectionTimeoutSeconds, err := strconv.Atoi(dbConfig.ConnectionTimeoutSeconds)

	if err == nil {
		db.SetConnMaxLifetime(time.Second * time.Duration(connectionTimeoutSeconds))
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	stats := db.Stats()
	logger.Info(fmt.Sprintf("{DUSER:%v, Idle:%v, OpenConnections:%v, InUse:%v, WaitCount:%v, WaitDuration:%v, MaxIdleClosed:%v, MaxLifetimeClosed:%v}",
		dbConfig.Username, stats.Idle, stats.OpenConnections, stats.InUse, stats.WaitCount, stats.WaitDuration, stats.MaxIdleClosed, stats.MaxLifetimeClosed))

	return db, nil
}

func ClosePostgresConnection(db *sql.DB) {
	db.Close()
}

func getConnectionString(dbConfig *structures.DbConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Db)
}
