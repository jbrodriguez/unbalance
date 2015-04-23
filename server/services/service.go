package services

import (
	"apertoire.net/unbalance/server/model"
)

type Service interface {
	ConfigChanged(conf *model.Config)
}
