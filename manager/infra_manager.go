package manager

import (
	"fmt"

	"calibration-system.com/config"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type InfraManager interface {
	Conn() *gorm.DB
	Migrate(model ...any) error
}

type infraManager struct {
	db  *gorm.DB
	cfg *config.Config
}

func (i *infraManager) Conn() *gorm.DB {
	return i.db
}

func (i *infraManager) initDb() error {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		i.cfg.Host, i.cfg.Port, i.cfg.User, i.cfg.DbConfig.Password, i.cfg.Name)
	db, err := gorm.Open(postgres.Open(psqlConn), &gorm.Config{})
	if err != nil {
		return err
	}
	i.db = db.Debug()
	return nil
}

func (i *infraManager) Migrate(model ...any) error {
	err := i.Conn().AutoMigrate(model...)
	if err != nil {
		return err
	}
	return nil

}

func NewInfraManager(cfg *config.Config) (InfraManager, error) {
	conn := &infraManager{
		cfg: cfg,
	}
	err := conn.initDb()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
