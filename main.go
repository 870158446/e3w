package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/soyking/e3w/conf"
	"github.com/soyking/e3w/e3ch"
	"github.com/soyking/e3w/routers"
	"go.etcd.io/etcd/version"
)

const (
	programName    = "e3w"
	programVersion = "0.1.0"
)

var configFilepath string

func init() {
	flag.StringVar(&configFilepath, "conf", "conf/config.default.ini", "config file path")
	rev := flag.Bool("rev", false, "print rev")
	flag.Parse()

	if *rev {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			programName, programVersion,
			version.Version,
		)
		os.Exit(0)
	}
}

func main() {
	config, err := conf.Init(configFilepath)
	if err != nil {
		panic(err)
	}

	client, err := e3ch.NewE3chClient(config)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.UseRawPath = true
	routers.InitRouters(router, config, client)
	router.Run(":" + config.Port)
}
