package config

import (
    "fmt"
    "os"

    "github.com/spf13/viper"
)

var V *viper.Viper

func ReadConfig(confFile string) *viper.Viper {
    file, e := os.Open(confFile)
    if e != nil {
        fmt.Println("ReadConfigFile Error:", e)
    }

    var v viper.Viper

    v.SetConfigType("prop")
    v.ReadConfig(file)

    return &v
}
