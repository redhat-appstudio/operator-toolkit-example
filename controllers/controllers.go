package controllers

import (
	"github.com/redhat-appstudio/operator-toolkit-example/controllers/bar"
	"github.com/redhat-appstudio/operator-toolkit-example/controllers/foo"
	"github.com/redhat-appstudio/operator-toolkit/controller"
)

// EnabledControllers is a slice containing references to all the controllers that have to be registered
var EnabledControllers = []controller.Controller{
	&bar.Controller{},
	&foo.Controller{},
}
