package discovery

import (
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"time"

	"github.com/promoboxx/go-cache/src/cache"
	"github.com/promoboxx/go-cache/src/storage/memory"
)

type Finder interface {
	FindService(name string) (string, error)
	FindHostPort(name string) (string, uint16, error)
}

type srvFinder struct {
	useTLS          bool
	hostSuffix      string
	proxyPort       uint16
	proxyPortString string
	c               cache.Cache
}

func NewSrvFinder(hostSuffix string, proxyPort uint16, useTLS bool) Finder {
	c := cache.New(memory.NewStorage(time.Second*5, time.Second*10), time.Second*5, true, nil)
	return &srvFinder{useTLS: useTLS, hostSuffix: hostSuffix, proxyPort: proxyPort, proxyPortString: fmt.Sprintf("%d", proxyPort), c: c}
}

func (f *srvFinder) FindService(name string) (string, error) {
	u := url.URL{}
	if f.useTLS {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	u.Host = net.JoinHostPort("proxy."+f.hostSuffix, f.proxyPortString)
	u.Path = name
	return u.String(), nil
}

func (f *srvFinder) FindHostPort(name string) (string, uint16, error) {
	var addrs []*net.SRV
	loader := func() (interface{}, error) {
		_, addrs, err := net.LookupSRV(name, "tcp", f.hostSuffix)
		return addrs, err
	}
	err := f.c.GetAndLoad(name, &addrs, loader)
	if err != nil {
		return "", 0, fmt.Errorf("Error during SRV lookup for (%s - %s): %v", name, f.hostSuffix, err)
	}
	if len(addrs) == 0 {
		return "", 0, fmt.Errorf("Error during SRV lookup for (%s - %s): no addresses found", name, f.hostSuffix)
	}

	srv := addrs[rand.Intn(len(addrs))]

	return srv.Target, srv.Port, nil
}
