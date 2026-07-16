package main

import (
	"server/core"
	"server/flag"
	"server/global"
	"server/initialize"
)

func main() {
	global.Configs = core.InitConf()
	global.Log = core.InitLogger()
	initialize.InitOther()
	global.DB = initialize.InitGorm()
	global.Redis = initialize.InitRedis()
	global.ESClient = initialize.InitEs()

	defer global.Redis.Close()

	flag.InitFlag()

	initialize.InitCron()

	core.RunServer()
}
