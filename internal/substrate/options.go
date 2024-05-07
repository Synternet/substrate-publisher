package substrate

import (
	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/service"
)

var RPCAPIParam = "rpc"

func WithRPCAPI(url string) options.Option {
	return func(o *options.Options) {
		service.WithParam(RPCAPIParam, url)(o)
	}
}

func (p *Publisher) RPCApi() string {
	return options.Param(p.Options, RPCAPIParam, "wss://rpc.polkadot.io")
}
