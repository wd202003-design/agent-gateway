package gateway

import (
 "encoding/json"
 "io"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestEndpoints(t *testing.T){
 for _,path:=range []string{"/v1/chat/completions","/v1/responses"}{
  for _,stream:=range []bool{false,true}{
   t.Run(path+map[bool]string{true:"/stream",false:"/json"}[stream],func(t *testing.T){
    upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
     if r.Header.Get("Authorization")!="Bearer upstream-secret" {t.Error("upstream authorization incorrect")}
     var v map[string]any
     if e:=json.NewDecoder(r.Body).Decode(&v);e!=nil || v["model"]!="native-model"{t.Errorf("upstream payload: %v %v",v,e)}
     if stream {w.Header().Set("Content-Type","text/event-stream");_,_=io.WriteString(w,"data: {\"ok\":true}\n\ndata: [DONE]\n\n")} else {w.Header().Set("Content-Type","application/json");_,_=io.WriteString(w,"{\"id\":\"test\"}")}
    }))
    defer upstream.Close()
    h:=New(Config{ClientKey:"client-secret",UpstreamKey:"upstream-secret",UpstreamURL:upstream.URL,PublicModel:"public-model",UpstreamModel:"native-model"})
    body:=`{"model":"public-model","stream":`+map[bool]string{true:"true",false:"false"}[stream]+`}`
    r:=httptest.NewRequest("POST",path,strings.NewReader(body));r.Header.Set("Authorization","Bearer client-secret");r.Header.Set("Content-Type","application/json")
    w:=httptest.NewRecorder();h.ServeHTTP(w,r)
    if w.Code!=200 {t.Fatalf("status=%d body=%s",w.Code,w.Body.String())}
    if stream && !strings.Contains(w.Body.String(),"[DONE]"){t.Fatal("missing SSE termination")}
   })
  }
 }
}
func TestAuthAndModel(t *testing.T){
 h:=New(Config{ClientKey:"secret",UpstreamKey:"up",UpstreamURL:"http://localhost",PublicModel:"pub",UpstreamModel:"native"})
 r:=httptest.NewRequest("POST","/v1/responses",strings.NewReader(`{"model":"pub"}`));r.Header.Set("Content-Type","application/json")
 w:=httptest.NewRecorder();h.ServeHTTP(w,r);if w.Code!=401{t.Fatal(w.Code)}
 r=httptest.NewRequest("POST","/v1/responses",strings.NewReader(`{"model":"wrong"}`));r.Header.Set("Authorization","Bearer secret");r.Header.Set("Content-Type","application/json")
 w=httptest.NewRecorder();h.ServeHTTP(w,r);if w.Code!=400{t.Fatal(w.Code)}
}
