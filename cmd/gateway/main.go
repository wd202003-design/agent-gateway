package main

import (
 "context"
 "log"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "time"

 "github.com/wd202003-design/agent-gateway/internal/gateway"
)

func main() {
 cfg, err := gateway.ConfigFromEnv()
 if err != nil { log.Fatal(err) }
 srv := &http.Server{Addr: cfg.ListenAddr, Handler: gateway.New(cfg), ReadHeaderTimeout: 10*time.Second}
 stop, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
 defer cancel()
 go func(){ <-stop.Done(); ctx, done := context.WithTimeout(context.Background(), 10*time.Second); defer done(); _ = srv.Shutdown(ctx) }()
 log.Printf("gateway listening on %s (development mode; billing disabled)",cfg.ListenAddr)
 if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Print(err); os.Exit(1) }
}
