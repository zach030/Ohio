package main

import (
	"flag"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
)

var env string

func init() {
	flag.StringVar(&env, "env", "", "default env name")
}

func main(){
	flag.Parse()
	log.Info("ohio-web-server start")
	paladin.Init()

}
