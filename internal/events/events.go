package events

import (
 "context"
 "encoding/json"
 "errors"
 "time"
)

type UsageEvent struct {
 EventID string `json:"event_id"`
 RequestID string `json:"request_id"`
 AttemptID string `json:"attempt_id"`
 TenantID string `json:"tenant_id"`
 SchemaVersion int `json:"schema_version"`
 InputTokens int64 `json:"input_tokens"`
 OutputTokens int64 `json:"output_tokens"`
 OccurredAt time.Time `json:"occurred_at"`
}
func (e UsageEvent) Validate()error{
 if e.EventID=="" || e.RequestID=="" || e.AttemptID=="" || e.TenantID=="" || e.SchemaVersion!=1 || e.InputTokens<0 || e.OutputTokens<0{return errors.New("invalid usage event")}
 return nil
}
func (e UsageEvent) Encode()([]byte,error){if err:=e.Validate();err!=nil{return nil,err};return json.Marshal(e)}
type Publisher interface { Publish(ctx context.Context,subject string,eventID string,payload []byte)error }
func PublishUsage(ctx context.Context,p Publisher,e UsageEvent)error{
 payload,err:=e.Encode();if err!=nil{return err}
 return p.Publish(ctx,"gateway.usage.v1",e.EventID,payload)
}
// Production publishers must use NATS JetStream with explicit ACK and dedupe;
// callers must persist an outbox record before publishing.
