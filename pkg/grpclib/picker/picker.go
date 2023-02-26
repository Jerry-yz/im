package picker

import (
	"context"
	"errors"
	"learn-im/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const AddrPickerName = "addr"

type addrKey struct{}

var ErrNotSubConnSelect = errors.New("no sub conn select")

func init() {
	balancer.Register(newBuilder())
}

func ContextWithAddr(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, addrKey{}, addr)
}

type addrPickerBuilder struct {
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(AddrPickerName, &addrPickerBuilder{}, base.Config{HealthCheck: true})
}

func (a *addrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	subConns := make(map[string]balancer.SubConn, len(info.ReadySCs))
	for k, v := range info.ReadySCs {
		subConns[v.Address.Addr] = k
	}
	return &addrPicker{
		SubConns: subConns,
	}
}

type addrPicker struct {
	SubConns map[string]balancer.SubConn
}

func (p *addrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	pr := balancer.PickResult{}

	address := info.Ctx.Value(addrKey{}).(string)
	sc, ok := p.SubConns[address]
	if !ok {
		logger.Logger.Error("Pick error", zap.String("address", address), zap.Any("subConnes", p.SubConns))
		return pr, ErrNotSubConnSelect
	}
	pr.SubConn = sc
	return pr, nil
}
