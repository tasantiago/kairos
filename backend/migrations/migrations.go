package migrations

import (
 "context"
 "embed"
 "fmt"
 "sort"
 "github.com/jackc/pgx/v5/pgxpool"
)
//go:embed *.sql
var files embed.FS

func Apply(ctx context.Context, db *pgxpool.Pool) error {
 tx, err := db.Begin(ctx); if err != nil { return err }; defer tx.Rollback(ctx)
 if _,err=tx.Exec(ctx,"SELECT pg_advisory_xact_lock(19482630)");err!=nil{return err}
 if _,err=tx.Exec(ctx,"CREATE TABLE IF NOT EXISTS schema_migrations (nome text PRIMARY KEY, aplicado_em timestamptz NOT NULL DEFAULT now())");err!=nil{return err}
 entries,err:=files.ReadDir(".");if err!=nil{return err};sort.Slice(entries,func(i,j int)bool{return entries[i].Name()<entries[j].Name()})
 for _,e:=range entries{var exists bool;if err=tx.QueryRow(ctx,"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE nome=$1)",e.Name()).Scan(&exists);err!=nil{return err};if exists{continue};sql,err:=files.ReadFile(e.Name());if err!=nil{return err};if _,err=tx.Exec(ctx,string(sql));err!=nil{return fmt.Errorf("migration %s: %w",e.Name(),err)};if _,err=tx.Exec(ctx,"INSERT INTO schema_migrations(nome) VALUES($1)",e.Name());err!=nil{return err}}
 return tx.Commit(ctx)
}
