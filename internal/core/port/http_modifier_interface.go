package port

import (
    "net/http"

    "personal-api-gateway/internal/core/domain"
)

type HttpModifierPortInterface interface {
    ModifyRequestHeader(req *http.Request, endpointConfig domain.EndpointObject)
}
