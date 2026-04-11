package sqlite

import (
	"database/sql"
	"os"
	"path"

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

func New(config string) (driver.Storager, error) {
	debug := false

	file := path.Join(config, "data.db")

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	f.Close()

	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	db.LogMode(debug)

	db.AutoMigrate(&AppConfig{}, &User{}, &Group{}, &Resource{}, &Collection{}, &Policy{}, &ContainerOwner{})

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	c.DB.Close()
}
