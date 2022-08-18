package mockgen

import "github.com/splice/platform/services/echo/rpc/echoservice/v1"

func init() {
	registrar.AddService(&echoservice.EchoService{})
}
