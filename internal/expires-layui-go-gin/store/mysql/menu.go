package mysql

import (
	"fmt"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getMenuData() ([]model.Menu, error) {
	filePath := "configs/sql-data/menu.yml"
	menus := make([]model.Menu, 0)
	viper.SetConfigFile(filePath)
	viper.SetDefault("is_shown", 1)
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		zap.L().Error("database: viper.ReadInConfig failed")
		return nil, err
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.UnmarshalKey("menu", &menus); err != nil {
		zap.L().Error("database: viper.Unmarshal failed")
		return nil, err
	}
	zap.L().Info(fmt.Sprintf("database: menus data is %+v", menus))
	zap.L().Info(fmt.Sprintf("database: menus length is %d", len(menus)))
	return menus, nil
}
