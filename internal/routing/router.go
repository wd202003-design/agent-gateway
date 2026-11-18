package routing

import (
 "errors"
 "math"
 "sort"
 "sync"
 "time"
)

var ErrNoCandidate=errors.New("no compatible healthy deployment")
type Candidate struct {
 ID, PublicModel, Protocol, Region string
 Capabilities map[string]bool
 Weight float64
 CostPerMillion float64
 TTFTMillis float64
 ErrorRate float64
 Available bool
}
type Request struct {
 PublicModel, Protocol, Region string
 RequiredCapabilities []string
 AllowCrossModel bool
 AllowedSubstitutes map[string]bool
 MaxCostPerMillion float64
}
type Decision struct { Selected Candidate; Score float64; Considered int }
func Select(req Request, candidates []Candidate) (Decision,error){
 eligible:=make([]Decision,0,len(candidates))
 for _,c:=range candidates {
  if !c.Available || c.Protocol!=req.Protocol || (req.Region!="" && c.Region!=req.Region) {continue}
  if c.PublicModel!=req.PublicModel && (!req.AllowCrossModel || !req.AllowedSubstitutes[c.PublicModel]) {continue}
  if req.MaxCostPerMillion>0 && c.CostPerMillion>req.MaxCostPerMillion {continue}
  valid:=true
  for _,cap:=range req.RequiredCapabilities {if !c.Capabilities[cap]{valid=false;break}}
  if !valid {continue}
  score:=math.Max(0,c.Weight)*0.4 + (1-math.Min(1,math.Max(0,c.ErrorRate)))*40 - math.Max(0,c.TTFTMillis)*0.002 - math.Max(0,c.CostPerMillion)*0.1
  eligible=append(eligible,Decision{Selected:c,Score:score})
 }
 if len(eligible)==0{return Decision{},ErrNoCandidate}
 sort.SliceStable(eligible,func(i,j int)bool{if eligible[i].Score==eligible[j].Score{return eligible[i].Selected.ID<eligible[j].Selected.ID};return eligible[i].Score>eligible[j].Score})
 result:=eligible[0];result.Considered=len(eligible);return result,nil
}
type Breaker struct { mu sync.Mutex; failures int; openedAt time.Time; threshold int; cooldown time.Duration }
func NewBreaker(threshold int,cooldown time.Duration)*Breaker {if threshold<1{threshold=3};return &Breaker{threshold:threshold,cooldown:cooldown}}
func (b *Breaker) Allow(now time.Time)bool{b.mu.Lock();defer b.mu.Unlock();return b.openedAt.IsZero() || now.Sub(b.openedAt)>=b.cooldown}
func (b *Breaker) Report(success bool,now time.Time){b.mu.Lock();defer b.mu.Unlock();if success{b.failures=0;b.openedAt=time.Time{};return};b.failures++;if b.failures>=b.threshold{b.openedAt=now}}
