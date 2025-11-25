package domain

type JsonConfig struct {
    LogEndpointUrl bool             `json:"log_endpoint_url"`
    Endpoints      []EndpointObject `json:"endpoints"`
}

type EndpointObject struct {
    Endpoint       string   `json:"endpoint"`
    AllowedHeaders []string `json:"allowed_headers"`
}
