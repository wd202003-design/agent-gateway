package budget

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
)

var ErrExceeded=errors.New("budget exceeded")
type Repository struct{ DB *sql.DB }

// Reserve atomically checks the PostgreSQL authoritative balance.
// Redis coordination is a separate optimization and cannot override this transaction.
func (r Repository) Reserve(ctx context.Context,accountID,requestID,reservationID string,amount string)error{
 if r.DB==nil{return errors.New("postgres unavailable: admission denied")}
 tx,err:=r.DB.BeginTx(ctx,&sql.TxOptions{Isolation:sql.LevelSerializable})
 if err!=nil{return err}
 defer tx.Rollback()
 res,err:=tx.ExecContext(ctx,`UPDATE budget_accounts SET reserved_amount=reserved_amount+$1::numeric WHERE account_id=$2::uuid AND spent_amount+reserved_amount+$1::numeric<=limit_amount`,amount,accountID)
 if err!=nil{return err}
 n,err:=res.RowsAffected();if err!=nil{return err};if n!=1{return ErrExceeded}
 _,err=tx.ExecContext(ctx,`INSERT INTO budget_reservations(reservation_id,account_id,request_id,amount,status,expires_at) VALUES ($1::uuid,$2::uuid,$3::uuid,$4::numeric,'RESERVED',now()+interval '10 minutes')`,reservationID,accountID,requestID,amount)
 if err!=nil{return fmt.Errorf("reserve insert: %w",err)}
 return tx.Commit()
}

// Capture must only be called after usage/charge computation. It is idempotent
// by reservation status and the ledger idempotency key.
func (r Repository) Capture(ctx context.Context,reservationID,ledgerEntryID,idempotencyKey,charge string)error{
 if r.DB==nil{return errors.New("postgres unavailable")}
 tx,err:=r.DB.BeginTx(ctx,&sql.TxOptions{Isolation:sql.LevelSerializable});if err!=nil{return err};defer tx.Rollback()
 var accountID,requestID,status string
 var reserved string
 err=tx.QueryRowContext(ctx,`SELECT account_id::text,request_id::text,status,amount::text FROM budget_reservations WHERE reservation_id=$1::uuid FOR UPDATE`,reservationID).Scan(&accountID,&requestID,&status,&reserved)
 if err!=nil{return err}
 if status=="CAPTURED"{return nil}
 if status!="RESERVED"{return errors.New("invalid reservation state")}
 // Captured charge cannot exceed reservation in this MVP; upstream overage requires explicit reconciliation.
 var allowed bool
 err=tx.QueryRowContext(ctx,`SELECT $1::numeric<=$2::numeric`,charge,reserved).Scan(&allowed)
 if err!=nil{return err};if !allowed{return ErrExceeded}
 _,err=tx.ExecContext(ctx,`UPDATE budget_accounts SET reserved_amount=reserved_amount-$1::numeric,spent_amount=spent_amount+$2::numeric WHERE account_id=$3::uuid`,reserved,charge,accountID)
 if err!=nil{return err}
 _,err=tx.ExecContext(ctx,`UPDATE budget_reservations SET status='CAPTURED' WHERE reservation_id=$1::uuid`,reservationID)
 if err!=nil{return err}
 _,err=tx.ExecContext(ctx,`INSERT INTO billing_ledger(entry_id,account_id,request_id,amount,entry_type,idempotency_key) VALUES($1::uuid,$2::uuid,$3::uuid,$4::numeric,'CHARGE',$5)`,ledgerEntryID,accountID,requestID,charge,idempotencyKey)
 if err!=nil{return err}
 return tx.Commit()
}
