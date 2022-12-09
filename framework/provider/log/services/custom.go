package services

import (
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/contract"
	"io"
)

type HeroCustomLog struct {
	HeroLog
}

func NewHeroCustomLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)

	log := &HeroCustomLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(output)
	log.c = c
	return log, nil
}
