package sqlite

import (
	"database/sql"
	"os"
	"path"
	"sync"

	"github.com/jonasbroms/hbm/storage"
	"github.com/jonasbroms/hbm/storage/driver"

	"github.com/jinzhu/gorm"
	modsqlite "modernc.org/sqlite"
)

func init() {
	// modernc registers as "sqlite"; re-register under "sqlite3" for GORM v1 compatibility
	sql.Register("sqlite3", &modsqlite.Driver{})
	storage.RegisterDriver("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

var (
	dbInstances = map[string]*gorm.DB{}
	dbMu        sync.Mutex
)

// New returns a Storager backed by a shared *gorm.DB for the given path.
// The DB is opened once and reused across all callers; this avoids SQLITE_BUSY
// races that occur when concurrent requests each open their own connection.
func New(config string) (driver.Storager, error) {
	file := path.Join(config, "data.db")

	dbMu.Lock()
	defer dbMu.Unlock()

	if db, ok := dbInstances[file]; ok {
		return &Config{DB: db}, nil
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	f.Close()

	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxOpenConns(1)
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA busy_timeout=5000")
	db.LogMode(false)

	db.AutoMigrate(&AppConfig{}, &User{}, &Group{}, &Resource{}, &Collection{}, &Policy{}, &ContainerOwner{})

	dbInstances[file] = db
	return &Config{DB: db}, nil
}

// End is a no-op: the shared connection stays open for the lifetime of the process.
func (c *Config) End() {}
