package config

import (
	"SuhoCoin/util"
	"github.com/spf13/viper"
	"os"
)

var V *viper.Viper

func ReadConfig(confFile string) *viper.Viper {
	file, e := os.Open(confFile)
	err.ERR("ReadConfigFile Error:", e)
	var v viper.Viper
	v.SetConfigType("prop")
	v.ReadConfig(file)
	return &v
}
