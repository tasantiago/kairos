package utils

import (
 "encoding/json"
 "errors"
 "log/slog"
 "net/http"
 "io"
 "fluxo/services"
 "github.com/jackc/pgx/v5"
 "github.com/jackc/pgx/v5/pgconn"
)
func JSON(w http.ResponseWriter,status int,value any){w.Header().Set("Content-Type","application/json; charset=utf-8");w.Header().Set("Cache-Control","no-store");w.WriteHeader(status);if value!=nil{json.NewEncoder(w).Encode(value)}}
func Error(w http.ResponseWriter,status int,message string){JSON(w,status,map[string]string{"erro":message})}
func Decode(w http.ResponseWriter,r *http.Request,value any)bool{r.Body=http.MaxBytesReader(w,r.Body,1<<20);decoder:=json.NewDecoder(r.Body);if err:=decoder.Decode(value);err!=nil{Error(w,400,"Dados inválidos.");return false};if decoder.Decode(new(any))!=io.EOF{Error(w,400,"Envie um único objeto JSON.");return false};return true}
func DBError(w http.ResponseWriter,err error){var pg *pgconn.PgError;var editorial *services.WorkflowError;switch{case errors.As(err,&editorial):Error(w,editorial.Status,editorial.Message);case errors.Is(err,pgx.ErrNoRows):Error(w,404,"Registro não encontrado.");case errors.As(err,&pg)&&pg.Code=="23505":Error(w,409,"Este nome ou e-mail já está cadastrado.");case errors.As(err,&pg)&&pg.Code=="23503":Error(w,409,"Registro em uso ou referência inválida.");case errors.As(err,&pg)&&pg.Code=="22P02":Error(w,400,"Identificador inválido.");case err.Error()=="forbidden":Error(w,403,"Você não tem acesso a esta tarefa.");case err.Error()=="assignee_sector":Error(w,400,"O responsável e o revisor devem pertencer ao setor selecionado.");default:slog.Error("database operation failed","error",err);Error(w,500,"Não foi possível concluir a operação.")}}
