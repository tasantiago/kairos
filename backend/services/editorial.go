package services

import (
 "net/url"
 "strings"
 "time"
 "fluxo/models"
)

type WorkflowError struct { Message string; Status int }
func (e *WorkflowError) Error() string { return e.Message }
func workflowError(message string,status int) error { return &WorkflowError{Message:message,Status:status} }
var stageNames=map[string]string{"pauta":"Pauta","producao":"Produção","edicao":"Edição","revisao":"Revisão","aprovacao":"Aprovação","pronto_publicacao":"Pronto para publicar","publicado":"Publicado","encerrada":"Encerrada (legado)"}
var transitions=map[string][]string{"pauta":{"producao"},"producao":{"pauta","edicao","revisao"},"edicao":{"producao","revisao"},"revisao":{"producao","edicao","aprovacao"},"aprovacao":{"revisao","edicao","pronto_publicacao"},"pronto_publicacao":{"revisao","publicado"},"publicado":{"revisao"},"encerrada":{"pauta"}}
func includes(values []string,value string)bool{for _,v:=range values{if v==value{return true}};return false}
func StageName(stage string)string{return stageNames[stage]}
func ValidLink(value string)bool{if value==""{return true};if len(value)>2048{return false};u,err:=url.Parse(value);return err==nil&&(u.Scheme=="https"||u.Scheme=="http")&&u.Hostname()!=""&&u.User==nil}
func Deadline(t models.Tarefa)time.Time{if t.PrazoEtapa!=nil{return *t.PrazoEtapa};return t.Prazo}
func IsLate(t models.Tarefa,now time.Time)bool{return t.Status!="concluida"&&(Deadline(t).Before(now)||(t.PublicacaoPrevista!=nil&&t.PublicacaoPrevista.Before(now)))}
func StatusForStage(stage string)string{if stage=="pauta"{return "pendente"};if stage=="publicado"||stage=="encerrada"{return "concluida"};return "em_andamento"}

// ApplyEditorial owns transitions and approval metadata; client-supplied signoffs are ignored.
func ApplyEditorial(previous *models.Tarefa,t models.Tarefa,u models.Usuario,now time.Time)(models.Tarefa,error){
 if t.TipoConteudo==""{t.TipoConteudo="outro"};if t.Etapa==""{t.Etapa="pauta";if previous!=nil{t.Etapa=previous.Etapa}};if t.Canal==""{t.Canal="Comunicação interna"}
 if !includes([]string{"video","materia","arte","outro"},t.TipoConteudo)||stageNames[t.Etapa]==""||!includes([]string{"Portal institucional","Instagram","YouTube","Facebook","Imprensa","Comunicação interna"},t.Canal){return t,workflowError("Tipo, etapa ou canal editorial inválido.",400)}
 if !ValidLink(t.LinkMaterial)||!ValidLink(t.LinkPublicado){return t,workflowError("Use links HTTP ou HTTPS válidos, sem credenciais.",400)}
 if t.Prazo.IsZero()||(t.PrazoEtapa!=nil&&t.PrazoEtapa.After(t.Prazo))||(t.PublicacaoPrevista!=nil&&t.Prazo.After(*t.PublicacaoPrevista)){return t,workflowError("O prazo da etapa deve ser até a entrega; a entrega deve ocorrer até a publicação prevista.",400)}
 t.AprovadoEm=nil;t.AprovadorID=nil;t.PublicadoEm=nil
 if previous==nil{
  if t.Etapa!="pauta"{return t,workflowError("Novas demandas começam na pauta.",400)};t.Versao=1;t.LinkPublicado=""
 }else{
  if t.Versao!=previous.Versao{return t,workflowError("Esta demanda mudou. Reabra os detalhes antes de salvar.",409)}
  t.AprovadoEm=previous.AprovadoEm;t.AprovadorID=previous.AprovadorID;t.PublicadoEm=previous.PublicadoEm
  changed:=t.Titulo!=previous.Titulo||t.Descricao!=previous.Descricao||t.TipoConteudo!=previous.TipoConteudo||t.Canal!=previous.Canal||t.LinkMaterial!=previous.LinkMaterial||t.SetorID!=previous.SetorID
  if t.Etapa!=previous.Etapa&&!includes(transitions[previous.Etapa],t.Etapa){return t,workflowError("Siga as etapas de produção, revisão e aprovação antes de publicar.",400)}
  if previous.TipoConteudo=="video"&&previous.Etapa=="producao"&&t.Etapa=="revisao"{return t,workflowError("O vídeo precisa passar pela edição antes da revisão.",400)}
  if previous.Etapa=="revisao"&&t.Etapa=="aprovacao"&&previous.RevisorID!=nil&&*previous.RevisorID!=u.ID&&u.Perfil!="gestor"{return t,workflowError("A revisão deve ser concluída pelo revisor responsável ou gestor.",403)}
  final:=previous.Etapa=="publicado"||previous.Etapa=="encerrada"
  if final&&(t.Etapa!=previous.Etapa||changed)&&u.Perfil!="gestor"{return t,workflowError("Somente o gestor pode reabrir uma entrega finalizada.",403)}
  if ((previous.Etapa=="aprovacao"&&(t.Etapa=="revisao"||t.Etapa=="edicao"))||(final&&(t.Etapa!=previous.Etapa||changed)))&&strings.TrimSpace(t.Anotacao)==""{return t,workflowError("Registre o motivo dos ajustes ou da reabertura no histórico.",400)}
  if previous.Etapa=="aprovacao"&&t.Etapa=="pronto_publicacao"{
   if u.Perfil!="gestor"||u.SetorID!=previous.SetorID{return t,workflowError("A aprovação final é feita por um gestor do setor responsável.",403)}
   if changed{return t,workflowError("O conteúdo mudou durante a aprovação. Envie a nova versão para revisão.",400)}
   t.AprovadoEm=&now;t.AprovadorID=&u.ID
  }
  if changed&&(previous.AprovadoEm!=nil||previous.Etapa=="aprovacao"){
   if t.Etapa=="publicado"&&previous.Etapa!="publicado"{return t,workflowError("O conteúdo alterado precisa de nova revisão e aprovação.",400)}
   t.Etapa="revisao"
  }
  if t.Etapa!="pronto_publicacao"&&t.Etapa!="publicado"{t.AprovadoEm=nil;t.AprovadorID=nil;t.PublicadoEm=nil;t.LinkPublicado=""}
  if t.Etapa=="publicado"&&previous.Etapa!="publicado"{
   if previous.AprovadoEm==nil||previous.AprovadorID==nil||changed{return t,workflowError("A publicação exige uma aprovação válida da versão atual.",400)}
   if strings.TrimSpace(t.LinkPublicado)==""{return t,workflowError("Informe o link da publicação realizada.",400)}
   t.PublicadoEm=&now
  }
  if t.Etapa=="publicado"&&strings.TrimSpace(t.LinkPublicado)==""{return t,workflowError("Mantenha o link da publicação realizada.",400)}
  t.Versao=previous.Versao+1
 }
 t.Status=StatusForStage(t.Etapa);t.ConcluidoEm=nil
 if t.Status=="concluida"{t.ConcluidoEm=&now;if previous!=nil&&previous.ConcluidoEm!=nil{t.ConcluidoEm=previous.ConcluidoEm}}
 return t,nil
}
