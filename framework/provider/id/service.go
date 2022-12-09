package id

import (
	"github.com/rs/xid"
)

type HeroIDService struct {
}

func NewHeroIDService(params ...interface{}) (interface{}, error) {
	return &HeroIDService{}, nil
}

func (s *HeroIDService) NewID() string {
	return xid.New().String()
}
