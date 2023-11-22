package mysql

import (
	"errors"
	"fmt"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/pkg/config"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os/exec"
	"strings"
	"time"
)

var DB *gorm.DB

// InitDB 初始化mysql数据库
func InitDB(cfg *config.Settings) error {
	var err error
	err = InitMysqlDocker(cfg)
	if err != nil {
		zap.L().Error("database: run mysql docker failed")
		return err
	}

	user := cfg.MysqlSettings.User
	password := cfg.MysqlSettings.Password
	host := cfg.MysqlSettings.Host
	port := cfg.MysqlSettings.Port
	database := cfg.MysqlSettings.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, password, host, port, database)
	zap.L().Info("database: dsn info is " + dsn)
	dialector := mysql.Open(dsn)
	dbConfig := &gorm.Config{
		SkipDefaultTransaction: true,
	}

	maxRetries := 6
	retryDelay := 10 * time.Second
	for i := 0; i < maxRetries; i++ {
		zap.L().Info("database: MySQL is waiting to connect.")
		DB, err = gorm.Open(dialector, dbConfig)
		if err == nil || i == maxRetries-1 {
			break
		}
		time.Sleep(retryDelay)
	}
	if err != nil {
		// 如果所有重试都失败，返回错误
		return err
	}

	db, _ := DB.DB()
	err = db.Ping()
	if err != nil {
		return err
	}
	maxIdleConns := cfg.MysqlSettings.MaxIdleConns
	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns)
	} else {
		db.SetMaxIdleConns(10)
	}

	maxOpenConns := cfg.MysqlSettings.MaxOpenConns
	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
	} else {
		db.SetMaxOpenConns(100)
	}
	db.SetConnMaxLifetime(time.Hour)
	err = autoMigrate()
	if err != nil {
		zap.L().Error("database: auto migrate database schema failed")
		return err
	}

	err = insertDefaultRecord()
	if err != nil {
		zap.L().Error("database: insert default record failed")
		return err
	}

	return nil
}

// InitMysqlDocker 拉起mysql容器
func InitMysqlDocker(cfg *config.Settings) error {
	containerName := "mysql-expires"
	if isContainerRunning(containerName) {
		zap.L().Info("database: MySQL container is already running.")
		return nil
	}
	composeFilePath := "deployments/docker-compose.yml"
	if isContainerExistsButNotRunning(containerName) {
		zap.L().Info("database: MySQL container exists but not running.")
		cmd := exec.Command("docker-compose", "-f", composeFilePath, "start")
		err := cmd.Run()
		if err != nil {
			zap.L().Error("database: cmd Run failed.", zap.Error(err))
		}
	} else {
		zap.L().Info("database: MySQL container is not exists.")
		cmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
		err := cmd.Run()
		if err != nil {
			zap.L().Error("database: cmd Run failed.", zap.Error(err))
		}
	}

	// 等待容器启动
	if err := waitForContainerRunning(containerName); err != nil {
		zap.L().Error("database: MySQL container didn't start within the expected time.", zap.Error(err))
		return err
	}
	return nil
}

func isContainerRunning(containerName string) bool {
	// 检查是否容器正在运行
	cmd := exec.Command("docker", "ps", "--filter", "name="+containerName, "--format", "{{.Names}}")
	output, _ := cmd.CombinedOutput()
	containers := string(output)
	return strings.Contains(containers, containerName)
}

func isContainerExistsButNotRunning(containerName string) bool {
	// 检查是否容器正在运行
	cmd := exec.Command("docker", "ps", "-a", "--filter", "name="+containerName, "--format", "{{.Names}}")
	output, _ := cmd.CombinedOutput()
	containers := string(output)
	return strings.Contains(containers, containerName)
}

func waitForContainerRunning(containerName string) error {
	maxRetries := 60
	interval := 10 * time.Second
	for i := 0; i < maxRetries; i++ {
		if isContainerRunning(containerName) {
			zap.L().Info("database: MySQL container is already running.")
			return nil
		}
		zap.L().Info("database: MySQL container is waiting to run.")
		time.Sleep(interval)
	}
	return errors.New("Container did not start within the expected time")
}

func autoMigrate() error {
	tables := []model.Table{
		model.Menu{},
	}
	for _, table := range tables {
		err := DB.AutoMigrate(table)
		if err != nil {
			zap.L().Error(fmt.Sprintf("database: table %s auto migrate failed", table.TableName()))
			return err
		}
	}
	return nil
}

func insertDefaultRecord() error {
	var err error
	err = insertMenuRecord()
	if err != nil {
		return err
	}
	return nil
}

func insertMenuRecord() error {
	menus, err := getMenuData()
	if err != nil {
		return err
	}
	for _, menu := range menus {
		var count int64
		err := DB.Model(&menu).Where(&model.Menu{Code: menu.Code}).Count(&count).Error
		if err != nil {
			zap.L().Error("database: check test record failed")
			return err
		}
		if count == 0 {
			if err := DB.Debug().Create(&menu).Error; err != nil {
				zap.L().Error(fmt.Sprintf("database: insert %s record failed", menu.Code))
				return err
			}
		}
	}
	return nil
}
