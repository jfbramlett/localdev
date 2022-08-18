package mockgen

import "github.com/splice/platform/services/subscription/rpc/subscriptionservice/v1"

func init() {
	registrar.AddService(&subscriptionservice.SubscriptionService{})
}
