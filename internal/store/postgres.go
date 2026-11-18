package store

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "errors"
 "time"
)

// Repository uses database/sql. A PostgreSQL driver must be registered by the application.
type Repository struct { DB *sql.DB }
type Identity struct { TenantID, KeyID, TenantType string }

func (s Repository) Authenticate(ctx context.Context, rawKey string) (Identity,error) {
 if s.DB==nil { return Identity{}, errors.New("postgres not configured") }
 hash:=sha256.Sum256([]byte(rawKey))
 var id Identity
 err:=s.DB.QueryRowContext(ctx,`SELECT k.tenant_id::text,k.key_id::text,t.tenant_type FROM api_keys k JOIN tenants t ON t.tenant_id=k.tenant_id WHERE k.verifier_hash=$1 AND k.status='ACTIVE' AND t.status='ACTIVE' AND (k.expires_at IS NULL OR k.expires_at>now())`,hash[:]).Scan(&id.TenantID,&id.KeyID,&id.TenantType)
 if err!=nil {return Identity{},err}
 return id,nil
}
type Deployment struct { ID, Endpoint, CredentialRef, NativeModel, Protocol string }
func (s Repository) Deployments(ctx context.Context, nativeModel, protocol string) ([]Deployment,error){
 if s.DB==nil {return nil,errors.New("postgres not configured")}
 rows,err:=s.DB.QueryContext(ctx,`SELECT d.deployment_id::text,c.endpoint,c.credential_ref,m.native_name,d.protocol_family FROM channel_deployments d JOIN channels c ON c.channel_id=d.channel_id JOIN upstream_models m ON m.upstream_model_id=d.upstream_model_id WHERE m.native_name=$1 AND d.protocol_family=$2 AND d.status='ACTIVE' AND c.status='ACTIVE'`,nativeModel,protocol)
 if err!=nil{return nil,err};defer rows.Close()
 var out []Deployment
 for rows.Next(){var d Deployment;if err=rows.Scan(&d.ID,&d.Endpoint,&d.CredentialRef,&d.NativeModel,&d.Protocol);err!=nil{return nil,err};out=append(out,d)}
 return out,rows.Err()
}
func (s Repository) Health(ctx context.Context) error {
 if s.DB==nil {return errors.New("postgres not configured")}
 c,cancel:=context.WithTimeout(ctx,2*time.Second);defer cancel()
 return s.DB.PingContext(c)
}
