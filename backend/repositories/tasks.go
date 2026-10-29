package repositories

import (
 "context"
 "errors"
 "time"
 "fluxo/models"
 "fluxo/services"
 "github.com/jackc/pgx/v5"
 "github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct { DB *pgxpool.Pool }
const TaskColumns="id::text,codigo,titulo,descricao,origem,status,prioridade,servidor_id::text,setor_id::text,prazo,criado_em,atualizado_em,concluido_em,tipo_conteudo,etapa,canal,revisor_id::text,prazo_etapa,publicacao_prevista,link_material,link_publicado,aprovador_id::text,aprovado_em,publicado_em,versao"
func ScanTask(row pgx.Row)(models.Tarefa,error){var t models.Tarefa;err:=row.Scan(&t.ID,&t.Codigo,&t.Titulo,&t.Descricao,&t.Origem,&t.Status,&t.Prioridade,&t.ServidorID,&t.SetorID,&t.Prazo,&t.CriadoEm,&t.AtualizadoEm,&t.ConcluidoEm,&t.TipoConteudo,&t.Etapa,&t.Canal,&t.RevisorID,&t.PrazoEtapa,&t.PublicacaoPrevista,&t.LinkMaterial,&t.LinkPublicado,&t.AprovadorID,&t.AprovadoEm,&t.PublicadoEm,&t.Versao);return t,err}
func(r *Repository)Task(ctx context.Context,id string)(models.Tarefa,error){return ScanTask(r.DB.QueryRow(ctx,"SELECT "+TaskColumns+" FROM tarefas WHERE id=$1",id))}
func(r *Repository)Tasks(ctx context.Context,u models.Usuario)([]models.Tarefa,error){rows,err:=r.DB.Query(ctx,"SELECT "+TaskColumns+" FROM tarefas WHERE ($1::boolean OR setor_id=$2::uuid OR servidor_id=$3::uuid OR revisor_id=$3::uuid) ORDER BY criado_em DESC",u.Perfil=="gestor",u.SetorID,u.ID);if err!=nil{return nil,err};defer rows.Close();result:=[]models.Tarefa{};for rows.Next(){t,err:=ScanTask(rows);if err!=nil{return nil,err};result=append(result,t)};return result,rows.Err()}
func(r *Repository)SaveTask(ctx context.Context,t models.Tarefa,u models.Usuario,create bool)(models.Tarefa,error){
 note:=t.Anotacao
 tx,err:=r.DB.Begin(ctx);if err!=nil{return t,err};defer tx.Rollback(ctx)
 var previous models.Tarefa;var prior *models.Tarefa
 if !create{
  previous,err=ScanTask(tx.QueryRow(ctx,"SELECT "+TaskColumns+" FROM tarefas WHERE id=$1 FOR UPDATE",t.ID));if err!=nil{return t,err}
  if u.Perfil!="gestor"&&previous.SetorID!=u.SetorID&&(previous.ServidorID==nil||*previous.ServidorID!=u.ID)&&(previous.RevisorID==nil||*previous.RevisorID!=u.ID){return t,errors.New("forbidden")}
  prior=&previous
 }
 t,err=services.ApplyEditorial(prior,t,u,time.Now());if err!=nil{return t,err}
 for _,person:=range []*string{t.ServidorID,t.RevisorID}{
  if person==nil{continue};var sector string
  err=tx.QueryRow(ctx,"SELECT setor_id::text FROM usuarios WHERE id=$1 FOR SHARE",*person).Scan(&sector);if err!=nil{return t,err};if sector!=t.SetorID{return t,errors.New("assignee_sector")}
 }
 movement:="Pauta criada."
 if !create{movement="Dados editoriais atualizados.";if t.Etapa!=previous.Etapa{movement=services.StageName(previous.Etapa)+" → "+services.StageName(t.Etapa)+"."};if previous.AprovadoEm!=nil&&t.AprovadoEm==nil{movement+=" Aprovação invalidada; nova revisão necessária."};if t.Etapa=="publicado"&&previous.Etapa!="publicado"{movement+=" "+t.LinkPublicado}}
 if note!=""{movement+=" "+note}
 if create{
  t,err=ScanTask(tx.QueryRow(ctx,"INSERT INTO tarefas(titulo,descricao,origem,prioridade,servidor_id,setor_id,prazo,tipo_conteudo,etapa,canal,revisor_id,prazo_etapa,publicacao_prevista,link_material) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "+TaskColumns,t.Titulo,t.Descricao,t.Origem,t.Prioridade,t.ServidorID,t.SetorID,t.Prazo,t.TipoConteudo,t.Etapa,t.Canal,t.RevisorID,t.PrazoEtapa,t.PublicacaoPrevista,t.LinkMaterial))
 }else{
  t,err=ScanTask(tx.QueryRow(ctx,"UPDATE tarefas SET titulo=$2,descricao=$3,origem=$4,status=$5,prioridade=$6,servidor_id=$7,setor_id=$8,prazo=$9,atualizado_em=now(),concluido_em=$10,tipo_conteudo=$11,etapa=$12,canal=$13,revisor_id=$14,prazo_etapa=$15,publicacao_prevista=$16,link_material=$17,link_publicado=$18,aprovador_id=$19,aprovado_em=$20,publicado_em=$21,versao=$22 WHERE id=$1 RETURNING "+TaskColumns,t.ID,t.Titulo,t.Descricao,t.Origem,t.Status,t.Prioridade,t.ServidorID,t.SetorID,t.Prazo,t.ConcluidoEm,t.TipoConteudo,t.Etapa,t.Canal,t.RevisorID,t.PrazoEtapa,t.PublicacaoPrevista,t.LinkMaterial,t.LinkPublicado,t.AprovadorID,t.AprovadoEm,t.PublicadoEm,t.Versao))
 }
 if err!=nil{return t,err}
 var from *string;if !create{from=&previous.SetorID}
 _,err=tx.Exec(ctx,"INSERT INTO historicos_tarefa(tarefa_id,de_setor,para_setor,usuario_id,anotacao,etapa,prazo_etapa,responsavel_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8)",t.ID,from,t.SetorID,u.ID,movement,t.Etapa,services.Deadline(t),t.ServidorID);if err!=nil{return t,err}
 return t,tx.Commit(ctx)
}
