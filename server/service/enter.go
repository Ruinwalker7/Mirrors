package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/attack_awareness"
	"github.com/flipped-aurora/gin-vue-admin/server/service/environment"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup           system.ServiceGroup
	ExampleServiceGroup          example.ServiceGroup
	Attack_awarenessServiceGroup attack_awareness.ServiceGroup
	EnvironmentServiceGroup      environment.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
