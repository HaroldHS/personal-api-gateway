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

// @Summary: Modify headers of a HTTP Request object.
// @Description: Modify header of HTTP Request object before sending to Reverse Proxy object.
//
// @Param req   *http.Request
// @Param value domain.EndpointObject "e.g. {Scheme `http`, Host `localhost`, AllowedHeaders [`Host`, `Accept`]}"
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

