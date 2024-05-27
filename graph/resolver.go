package graph

import (
	"github.com/3XBAT/coments-system/pkg/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service *service.Service
}


func NewResolver(service *service.Service) *Resolver {
	return &Resolver{
		service: service,
	}
}
