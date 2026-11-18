package routing

import "testing"

func TestCrossModelDeniedByDefault(t *testing.T){
 c:=[]Candidate{{ID:"other",PublicModel:"replacement",Protocol:"responses",Available:true,Weight:10}}
 _,err:=Select(Request{PublicModel:"original",Protocol:"responses"},c)
 if err!=ErrNoCandidate{t.Fatal(err)}
 d,err:=Select(Request{PublicModel:"original",Protocol:"responses",AllowCrossModel:true,AllowedSubstitutes:map[string]bool{"replacement":true}},c)
 if err!=nil || d.Selected.ID!="other"{t.Fatalf("%+v %v",d,err)}
}
func TestCapabilityAndCostFilter(t *testing.T){
 c:=[]Candidate{{ID:"a",PublicModel:"p",Protocol:"chat",Available:true,Capabilities:map[string]bool{"tools":true},CostPerMillion:3}}
 _,err:=Select(Request{PublicModel:"p",Protocol:"chat",RequiredCapabilities:[]string{"tools"},MaxCostPerMillion:2},c)
 if err!=ErrNoCandidate{t.Fatal(err)}
}
