package gateway

import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestRejectInvalidRequests(t *testing.T) {
 h:=New(Config{ClientKey:"secret",UpstreamKey:"up",UpstreamURL:"http://127.0.0.1:9",PublicModel:"pub",UpstreamModel:"native"})
 cases:=[]struct{body,auth,contentType string;status int}{
  {`{"model":"pub"}`,"","application/json",401},
  {`{"model":"pub"}`,"Bearer secret","text/plain",415},
  {`{`,"Bearer secret","application/json",400},
  {`{"model":"other"}`,"Bearer secret","application/json",400},
 }
 for _,tc:=range cases{
  r:=httptest.NewRequest(http.MethodPost,"/v1/chat/completions",strings.NewReader(tc.body))
  r.Header.Set("Authorization",tc.auth);r.Header.Set("Content-Type",tc.contentType)
  w:=httptest.NewRecorder();h.ServeHTTP(w,r)
  if w.Code!=tc.status{t.Errorf("got %d want %d",w.Code,tc.status)}
 }
}

func TestCancellation(t *testing.T){
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){<-r.Context().Done()}))
 defer upstream.Close()
 h:=New(Config{ClientKey:"secret",UpstreamKey:"up",UpstreamURL:upstream.URL,PublicModel:"pub",UpstreamModel:"native"})
 ctx,cancel:=context.WithCancel(context.Background());cancel()
 r:=httptest.NewRequest("POST","/v1/responses",strings.NewReader(`{"model":"pub"}`)).WithContext(ctx)
 r.Header.Set("Authorization","Bearer secret");r.Header.Set("Content-Type","application/json")
 w:=httptest.NewRecorder();h.ServeHTTP(w,r)
 if w.Code!=502 && w.Code!=504 {t.Fatalf("unexpected code %d",w.Code)}
}
