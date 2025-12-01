package util

import (
    "errors"
    "net/http/httputil"
    "net/url"

    "personal-api-gateway/internal/core/domain"
)

type BasicHttpProxies struct {
    Proxies map[string]*httputil.ReverseProxy
}

func NewBasicHttpProxies(cfg domain.JsonConfig) (*BasicHttpProxies, error) {
    mappedProxies := make(map[string]*httputil.ReverseProxy, len(cfg.Endpoints))

    for key, obj := range cfg.Endpoints {
        targetUrl, err := url.Parse(obj.DestinationScheme + "://" + obj.DestinationHost + obj.DestinationPath)
        if err != nil {
            return &BasicHttpProxies{}, err
        }

        mappedProxies[key] = httputil.NewSingleHostReverseProxy(targetUrl)
    }

    return &BasicHttpProxies{
        Proxies: mappedProxies,
    }, nil
}

func (bhp *BasicHttpProxies) GetBasicHttpProxy(target string) (*httputil.ReverseProxy, error) {
    if httpProxy, ok := bhp.Proxies[target]; ok {
        return httpProxy, nil
    }

    return &httputil.ReverseProxy{}, errors.New("target proxy not found")
}

