package orm

import (
	"context"
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
	"time"
)

type HeroGorm struct {
	container framework.Container
	dbs       map[string]*gorm.DB //key为dsn，value是gorm.DB连接池

	lock *sync.RWMutex // 由于对dbs的操作时读多写少，最好使用读写锁
}

func NewHeroGorm(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}
	return HeroGorm{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil
}

// GetDB 获取DB实例
func (hg *HeroGorm) GetDB(option ...contract.DBOption) (*gorm.DB, error) {

	logger := hg.container.MustMake(contract.LogKey).(contract.Log)

	// 读取默认配置
	config := GetBaseConfig(hg.container)

	logService := hg.container.MustMake(contract.LogKey).(contract.Log)

	// 设置logger
	ormLogger := NewOrmLogger(logService)
	config.Config = &gorm.Config{
		Logger: ormLogger,
	}

	for _, opt := range option {
		if err := opt(hg.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn,就生成dsn
	if config.Dsn == "" {
		dsn, err := config.FormatDSN()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	//判断是否已经实例化了db
	hg.lock.RLock()
	if db, ok := hg.dbs[config.Dsn]; ok {
		hg.lock.RUnlock()
		return db, nil
	}
	hg.lock.RUnlock()

	//没有实例化，就进行实例化步骤
	hg.lock.Lock()
	defer hg.lock.Unlock()

	// 实例化gorm.DB
	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config)
	}

	// 设置对应的连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}

	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}

	if config.ConnMaxLifetime != "" {
		liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max lift time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(liftTime)
		}
	}

	if config.ConnMaxIdletime != "" {
		idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
		if err != nil {
			logger.Error(context.Background(), "conn max idle time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxIdleTime(idleTime)
		}
	}

	// 挂载到map中
	if err != nil {
		hg.dbs[config.Dsn] = db
	}
	return db, err
}
