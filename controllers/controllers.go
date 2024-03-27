package controllers

import (
	"github.com/konflux-ci/operator-toolkit-example/controllers/bar"
	"github.com/konflux-ci/operator-toolkit-example/controllers/foo"
	"github.com/konflux-ci/operator-toolkit/controller"
)

// EnabledControllers is a slice containing references to all the controllers that have to be registered
var EnabledControllers = []controller.Controller{
	&bar.Controller{},
	&foo.Controller{},
}
