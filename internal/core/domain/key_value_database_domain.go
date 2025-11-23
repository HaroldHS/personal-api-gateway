package domain

type KeyValueObject struct {
    ReqId string `json:"req_id"`
    Key   string `json:"key"`
    Value any    `json:"value"`
    Ttl   int64  `json:"ttl"`
}
