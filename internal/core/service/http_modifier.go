package service

import (
    "net/http"

    "personal-api-gateway/internal/core/domain"
)

type HttpModifierService struct {
}

func NewHttpModifier() *HttpModifierService {
    return &HttpModifierService{}
}

func (hms *HttpModifierService) ModifyRequestHeader(req *http.Request, endpointConfig domain.EndpointObject) {
    req.Host = endpointConfig.DestinationHost
    req.URL.Path = endpointConfig.DestinationPath
    req.RequestURI = endpointConfig.DestinationPath
    for header := range req.Header {
        notWhiteListed := true
        for _, allowedHeader := range endpointConfig.AllowedHeaders {
            if header == allowedHeader {
                notWhiteListed = false
            }
        }
        if notWhiteListed {
            req.Header.Del(header)
        }
    }
    req.Header.Set("X-Forwarded-For", req.RemoteAddr)
}

