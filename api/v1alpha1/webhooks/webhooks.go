package webhooks

import (
	"github.com/redhat-appstudio/operator-toolkit-example/api/v1alpha1/webhooks/bar"
	"github.com/redhat-appstudio/operator-toolkit/webhook"
)

// EnabledWebhooks is a slice containing references to all the webhooks that have to be registered
var EnabledWebhooks = []webhook.Webhook{
	&bar.Webhook{},
}
