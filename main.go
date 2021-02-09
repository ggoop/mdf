package main

import (
	"github.com/ggoop/mdf/bootstrap"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/utils"
)

func main() {
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}
	bootstrap.StartHttp()
}
