package configs

import (
	"github.com/ggoop/mdf/utils"
)

var Default *utils.Config

func init() {
	Default = utils.DefaultConfig
}
