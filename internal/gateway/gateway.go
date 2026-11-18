package gateway

import (
 "bytes"
 "context"
 "crypto/sha256"
 "crypto/subtle"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "log"
 "net/http"
 "net/url"
 "os"
 "strings"
 "time"
)

const maxBody = 4 << 20

// Config is intentionally a single-upstream development slice. Production billing is disabled.
type Config struct {
 ListenAddr string
 ClientKey string
 UpstreamKey string
 UpstreamURL string
 PublicModel string
 UpstreamModel string
 Timeout time.Duration
}

func ConfigFromEnv() (Config,error) {
 c:=Config{ListenAddr: getenv("GATEWAY_LISTEN",":8080"), ClientKey:os.Getenv("GATEWAY_CLIENT_KEY"), UpstreamKey:os.Getenv("UPSTREAM_API_KEY"), UpstreamURL:os.Getenv("UPSTREAM_BASE_URL"), PublicModel:getenv("PUBLIC_MODEL","gateway-default"), UpstreamModel:os.Getenv("UPSTREAM_MODEL"), Timeout:120*time.Second}
 if c.ClientKey=="" || c.UpstreamKey=="" || c.UpstreamModel=="" { return c,errors.New("GATEWAY_CLIENT_KEY, UPSTREAM_API_KEY and UPSTREAM_MODEL are required") }
 u,e:=url.Parse(c.UpstreamURL); if e!=nil || u==nil || (u.Scheme!="http" && u.Scheme!="https") || u.Host=="" || u.User!=nil || u.RawQuery!="" || u.Fragment!="" { return c,errors.New("UPSTREAM_BASE_URL must be an absolute http(s) URL without credentials, query or fragment") }
 if c.PublicModel=="" { return c,errors.New("PUBLIC_MODEL must not be empty") }
 return c,nil
}
func getenv(k,d string)string { if v:=os.Getenv(k);v!="" {return v};return d }

type Handler struct { cfg Config; client *http.Client }
func New(c Config) http.Handler {
 h:=&Handler{cfg:c,client:&http.Client{Timeout:c.Timeout,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}}
 mux:=http.NewServeMux()
 mux.HandleFunc("/healthz",func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(http.StatusOK);_,_=io.WriteString(w,"ok")})
 mux.HandleFunc("/readyz",func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(http.StatusOK);_,_=io.WriteString(w,"development-only")})
 mux.HandleFunc("/v1/models",h.models)
 mux.HandleFunc("/v1/chat/completions",h.proxy)
 mux.HandleFunc("/v1/responses",h.proxy)
 return mux
}
func (h *Handler) authorized(r *http.Request) bool {
 a:=r.Header.Get("Authorization")
 if !strings.HasPrefix(a,"Bearer ") {return false}
 given:=strings.TrimPrefix(a,"Bearer ")
 x,y:=sha256.Sum256([]byte(given)),sha256.Sum256([]byte(h.cfg.ClientKey))
 return subtle.ConstantTimeCompare(x[:],y[:])==1
}
func jsonError(w http.ResponseWriter,status int,code,msg string) {
 w.Header().Set("Content-Type","application/json")
 w.WriteHeader(status)
 _=json.NewEncoder(w).Encode(map[string]any{"error":map[string]string{"type":code,"message":msg,"code":code}})
}
func (h *Handler) models(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet {jsonError(w,405,"method_not_allowed","GET required");return}
 if !h.authorized(r){jsonError(w,401,"invalid_api_key","invalid API key");return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(map[string]any{"object":"list","data":[]any{map[string]any{"id":h.cfg.PublicModel,"object":"model","owned_by":"agent-gateway"}}})
}
func (h *Handler) proxy(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {jsonError(w,405,"method_not_allowed","POST required");return}
 if !h.authorized(r){jsonError(w,401,"invalid_api_key","invalid API key");return}
 if !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")),"application/json") {jsonError(w,415,"unsupported_media_type","application/json required");return}
 r.Body=http.MaxBytesReader(w,r.Body,maxBody)
 body,err:=io.ReadAll(r.Body)
 if err!=nil {jsonError(w,413,"request_too_large","request exceeds size limit");return}
 var payload map[string]json.RawMessage
 if err=json.Unmarshal(body,&payload);err!=nil || payload==nil {jsonError(w,400,"invalid_json","expected JSON object");return}
 var model string
 if err=json.Unmarshal(payload["model"],&model);err!=nil || model!=h.cfg.PublicModel {jsonError(w,400,"model_not_found","model not available");return}
 payload["model"],_=json.Marshal(h.cfg.UpstreamModel)
 rewritten,err:=json.Marshal(payload);if err!=nil {jsonError(w,400,"invalid_json","cannot encode request");return}
 target,err:=url.Parse(strings.TrimRight(h.cfg.UpstreamURL,"/") + r.URL.Path)
 if err!=nil {jsonError(w,500,"config_error","invalid upstream configuration");return}
 req,err:=http.NewRequestWithContext(r.Context(),http.MethodPost,target.String(),bytes.NewReader(rewritten))
 if err!=nil {jsonError(w,500,"request_error","cannot construct request");return}
 req.Header.Set("Authorization","Bearer "+h.cfg.UpstreamKey)
 req.Header.Set("Content-Type","application/json")
 req.Header.Set("Accept",r.Header.Get("Accept"))
 start:=time.Now()
 resp,err:=h.client.Do(req)
 if err!=nil {
  status:=http.StatusBadGateway; if errors.Is(err,context.DeadlineExceeded){status=http.StatusGatewayTimeout}
  jsonError(w,status,"upstream_error","upstream request failed");return
 }
 defer resp.Body.Close()
 for _,k:=range []string{"Content-Type","Cache-Control","X-Request-Id"} {if v:=resp.Header.Get(k);v!="" {w.Header().Set(k,v)}}
 w.Header().Set("X-Content-Type-Options","nosniff")
 w.WriteHeader(resp.StatusCode)
 flusher,_:=w.(http.Flusher)
 buf:=make([]byte,16*1024)
 var nbytes int64
 for {
  n,e:=resp.Body.Read(buf)
  if n>0 {wn,we:=w.Write(buf[:n]);nbytes+=int64(wn);if flusher!=nil {flusher.Flush()};if we!=nil {break}}
  if e!=nil {if e!=io.EOF {log.Printf("upstream read interrupted: path=%s status=%d",r.URL.Path,resp.StatusCode)};break}
 }
 log.Printf("request path=%s status=%d duration_ms=%d bytes=%d",r.URL.Path,resp.StatusCode,time.Since(start).Milliseconds(),nbytes)
}
var _ = fmt.Sprintf
