package operator

import (
	"game/duck/ut2"
	"net/url"
	"strings"

	"google.golang.org/grpc/resolver"
)

type ResolverBuilder struct {
	appid2res *ut2.SyncMap[string, *Resolver]
	mgr       *AppManager
}

func NewResolverBuilder(mgr *AppManager) *ResolverBuilder {
	return &ResolverBuilder{
		mgr:       mgr,
		appid2res: ut2.NewSyncMap[string, *Resolver](),
	}
}

func (mb *ResolverBuilder) Scheme() string {
	return "grpc"
}

func (mb *ResolverBuilder) Update(appid, addr string) {
	res, ok := mb.appid2res.Load(appid)
	if !ok {
		return
	}

	res.ResolveNow(resolver.ResolveNowOptions{})
}

func (mb *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		URL: target.URL,
		cc:  cc,
		mb:  mb,
	}

	mb.appid2res.Store(target.URL.Host, r)

	r.ResolveNow(resolver.ResolveNowOptions{})

	return r, nil
}

type Resolver struct {
	URL url.URL
	cc  resolver.ClientConn
	mb  *ResolverBuilder
}

func (r *Resolver) Close() {
	r.mb.appid2res.Delete(r.URL.Host)
}

func (r *Resolver) ResolveNow(opt resolver.ResolveNowOptions) {
	r.cc.UpdateState(resolver.State{Addresses: r.build()})
}

func (r *Resolver) build() (ret []resolver.Address) {
	app := r.mb.mgr.GetApp(r.URL.Host)

	addr := strings.ReplaceAll(app.Address, "grpc://", "")

	ret = append(ret, resolver.Address{Addr: addr})

	return
}
