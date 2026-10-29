package models

import "time"

type Usuario struct {
 ID string `json:"id"`
 Nome string `json:"nome"`
 Email string `json:"email"`
 Perfil string `json:"perfil"`
 SetorID string `json:"setor_id"`
 Senha string `json:"senha,omitempty"`
}
type Setor struct {
 ID string `json:"id"`
 Nome string `json:"nome"`
 Capacidade int `json:"capacidade"`
 Cor string `json:"cor"`
}
type Tarefa struct {
 ID string `json:"id"`
 Codigo string `json:"codigo"`
 Titulo string `json:"titulo"`
 Descricao string `json:"descricao"`
 Origem string `json:"origem"`
 Status string `json:"status"`
 Prioridade string `json:"prioridade"`
 ServidorID *string `json:"servidor_id"`
 SetorID string `json:"setor_id"`
 Prazo time.Time `json:"prazo"`
 CriadoEm time.Time `json:"criado_em"`
 AtualizadoEm time.Time `json:"atualizado_em"`
 ConcluidoEm *time.Time `json:"concluido_em"`
 Anotacao string `json:"anotacao,omitempty"`
 TipoConteudo string `json:"tipo_conteudo"`
 Etapa string `json:"etapa"`
 Canal string `json:"canal"`
 RevisorID *string `json:"revisor_id"`
 PrazoEtapa *time.Time `json:"prazo_etapa"`
 PublicacaoPrevista *time.Time `json:"publicacao_prevista"`
 LinkMaterial string `json:"link_material"`
 LinkPublicado string `json:"link_publicado"`
 AprovadorID *string `json:"aprovador_id"`
 AprovadoEm *time.Time `json:"aprovado_em"`
 PublicadoEm *time.Time `json:"publicado_em"`
 Versao int `json:"versao"`
}
type Historico struct {
 ID string `json:"id"`
 TarefaID string `json:"tarefa_id"`
 DeSetor *string `json:"de_setor"`
 ParaSetor *string `json:"para_setor"`
 UsuarioID *string `json:"usuario_id"`
 Anotacao string `json:"anotacao"`
 Data time.Time `json:"data"`
 Etapa *string `json:"etapa"`
 PrazoEtapa *time.Time `json:"prazo_etapa"`
 ResponsavelID *string `json:"responsavel_id"`
}
type Mensagem struct {
 ID string `json:"id"`
 TarefaID string `json:"tarefa_id"`
 UsuarioID *string `json:"usuario_id"`
 Mensagem string `json:"mensagem"`
 Data time.Time `json:"data"`
}
