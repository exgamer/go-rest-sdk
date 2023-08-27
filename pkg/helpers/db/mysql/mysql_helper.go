package helpers

import (
	"database/sql"
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	"log"
	"time"
)

func OpenMysqlConnection(dbConfig *structures.DbConfig) (*sql.DB, error) {
	// Open up database connection.
	db, err := sql.Open("mysql", getConnectionString(dbConfig))

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(dbConfig.MaxPoolConnections)
	db.SetMaxIdleConns(dbConfig.MaxIdlePoolConnections)
	db.SetConnMaxLifetime(time.Second * time.Duration(dbConfig.ConnectionTimeoutSeconds))

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	stats := db.Stats()
	logger.Info(fmt.Sprintf("{DUSER:%v, Idle:%v, OpenConnections:%v, InUse:%v, WaitCount:%v, WaitDuration:%v, MaxIdleClosed:%v, MaxLifetimeClosed:%v}",
		dbConfig.Username, stats.Idle, stats.OpenConnections, stats.InUse, stats.WaitCount, stats.WaitDuration, stats.MaxIdleClosed, stats.MaxLifetimeClosed))

	return db, nil
}

func CloseMysqlConnection(db *sql.DB) {
	db.Close()
}

func getConnectionString(dbConfig *structures.DbConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Db)
}
