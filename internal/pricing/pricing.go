package pricing

import (
 "errors"
 "math/big"
)

// Money is a rational decimal quantity; rounding happens only at the invoice boundary.
// Do not use float64 for prices or charges.
type Rate struct { PerMillionInput, PerMillionOutput string }
type Usage struct { InputTokens, OutputTokens int64 }
func Calculate(rate Rate,usage Usage)(string,error){
 if usage.InputTokens<0 || usage.OutputTokens<0{return "",errors.New("negative usage")}
 a,ok:=new(big.Rat).SetString(rate.PerMillionInput);if !ok || a.Sign()<0{return "",errors.New("invalid input rate")}
 b,ok:=new(big.Rat).SetString(rate.PerMillionOutput);if !ok || b.Sign()<0{return "",errors.New("invalid output rate")}
 x:=new(big.Rat).Mul(a,new(big.Rat).SetInt64(usage.InputTokens))
 y:=new(big.Rat).Mul(b,new(big.Rat).SetInt64(usage.OutputTokens))
 total:=new(big.Rat).Quo(new(big.Rat).Add(x,y),big.NewRat(1000000,1))
 return total.FloatString(12),nil
}
