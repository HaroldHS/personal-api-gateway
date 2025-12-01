package domain

type JsonConfig struct {
    LogEndpointUrl bool                      `json:"log_endpoint_url"`
    Endpoints      map[string]EndpointObject `json:"endpoints"`
}

type EndpointObject struct {
    DestinationScheme   string   `json:"destination_scheme"`
    DestinationHost     string   `json:"destination_host"`
    DestinationPath     string   `json:"destination_path"`
    AllowedHeaders      []string `json:"allowed_headers"`
}
