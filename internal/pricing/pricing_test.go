package pricing

import "testing"

func TestCalculate(t *testing.T){
 got,err:=Calculate(Rate{PerMillionInput:"1.0",PerMillionOutput:"4.0"},Usage{InputTokens:2000,OutputTokens:1000})
 if err!=nil || got!="0.006000000000"{t.Fatalf("got=%s err=%v",got,err)}
 _,err=Calculate(Rate{PerMillionInput:"-1",PerMillionOutput:"4"},Usage{})
 if err==nil{t.Fatal("negative price accepted")}
}
