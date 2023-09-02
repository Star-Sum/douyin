package main

import (
	"douyin/Dao/RedisDao"
	"douyin/Servlet"
)

func main() {
	Servlet.DbServerStartup()
	Servlet.NetServerStartup()
	RedisDao.FlushDb()
}
