package postgres

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// dbWrite dbRead dbLogger: variable for database
var (
	dbWrite, dbRead *gorm.DB
	dbLogger        *log.Logger
	isDebug         bool
	dbReadMu        sync.Mutex
)

// DBLogFormatter database log formatter
type DBLogFormatter struct {
	EnableColor bool
}

// Format function to format the database log
// entry log.Entry
func (f *DBLogFormatter) Format(entry *log.Entry) ([]byte, error) {
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	b := &bytes.Buffer{}
	if entry.Message != "" {
		m := entry.Message
		if f.EnableColor == false {
			for _, v := range []string{
				"\033[33m",
				"\033[35m",
				"\033[36;1m",
				"\033[31;1m",
				"\033[0m",
			} {
				m = strings.Replace(m, v, "", -1)
			}
		}
		b.WriteString(m)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

// InitDB function to initialize database log
func InitDB() {
	isDebug = false
	if os.Getenv("APP_DEBUG") == "1" {
		isDebug = true
	}
	fmt.Println(fmt.Sprintf("debug: %v", isDebug))

	if isDebug {
		dbLogger = log.New()
		dbLogger.Formatter = &DBLogFormatter{EnableColor: false}
	}
}

// GetWriteDB function to get writing access to database
func GetWriteDB() *gorm.DB {
	if dbWrite == nil {
		dbWrite = CreateDBConnection(fmt.Sprintf("host=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_DB_WRITE_HOST"), os.Getenv("POSTGRES_DB_WRITE_USER"), os.Getenv("POSTGRES_DB_WRITE_PASSWORD"), os.Getenv("POSTGRES_DB_WRITE_NAME")))
	}
	return dbWrite
}

// GetReadDB function to get reading access to database
func GetReadDB() *gorm.DB {
	dbReadMu.Lock()
	defer dbReadMu.Unlock()

	if dbRead == nil {
		//dbRead = CreateDBConnection(os.Getenv("DB_READ"))
		dbRead = CreateDBConnection(fmt.Sprintf("host=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_DB_READ_HOST"), os.Getenv("POSTGRES_DB_READ_USER"), os.Getenv("POSTGRES_DB_READ_PASSWORD"), os.Getenv("POSTGRES_DB_READ_NAME")))

	}
	return dbRead
}

// CreateDBConnection function to create database connection
func CreateDBConnection(descriptor string) *gorm.DB {
	db, err := gorm.Open("postgres", descriptor)
	if err != nil {
		defer db.Close()
		return db
	}

	maxOpenCons, _ := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONS"))

	// set max idle connection to zero
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(maxOpenCons)

	// set database log into file
	if isDebug {
		db.LogMode(true)
		db.SetLogger(gorm.Logger{dbLogger})
	}

	return db
}

// CloseDb function for closing database connection
func CloseDb() {
	if dbRead != nil {
		dbRead.Close()
		dbRead = nil
	}
	if dbWrite != nil {
		dbWrite.Close()
		dbWrite = nil
	}
}
