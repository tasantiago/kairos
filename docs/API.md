# API do Fluxo

Base `/api`. JSON UTF-8. Autenticação: `Authorization: Bearer <token>`. Erros: `{ "erro": "mensagem em português" }`. POST retorna 201, GET/PUT 200, DELETE 204. UUID inválido retorna 400; registro ausente, 404; conflito de dependência ou unicidade, 409. Listas vazias retornam `[]`.

| Método | Rota | Acesso e comportamento |
|---|---|---|
| GET | `/health` | Público; confirma acesso ao banco, 503 se indisponível (fora de `/api`) |
| POST | `/login` | Público; `{email, senha}` → `{token, usuario}` |
| GET | `/me` | Usuário autenticado |
| GET/POST | `/tarefas` | Listar demandas visíveis / criar no próprio setor; gestor pode criar em qualquer setor |
| GET/PUT/DELETE | `/tarefas/{id}` | Ler/editar se houver acesso; exclusão apenas por gestor |
| GET | `/tarefas/prioritaria` | Tarefa aberta com maior pontuação; `null` se não houver |
| GET | `/tarefas/{id}/historico` | Histórico cronológico, somente se houver acesso à tarefa |
| GET/POST | `/tarefas/{id}/chat` | Listar/enviar mensagem: `{mensagem}`; até 4000 caracteres |
| GET/POST | `/tarefas/{id}/anexos` | Listar/enviar multipart, campo `arquivo`, até 10 MB |
| GET/DELETE | `/anexos/{id}` | Download autenticado / excluir como gestor |
| GET/POST | `/usuarios` | Listar dados sem hash/senha / criar como gestor |
| GET/PUT/DELETE | `/usuarios/{id}` | Ler / editar ou excluir como gestor |
| GET/POST | `/setores` | Listar / criar como gestor |
| GET/PUT/DELETE | `/setores/{id}` | Ler / editar ou excluir como gestor |
| GET/PUT | `/config` | Ler pontuação / ajustar como gestor |
| GET | `/metricas` | Gestor; totais, carga, média, mapa de criação e títulos recorrentes |
| GET | `/relatorios` | Gestor; CSV UTF-8 com BOM, registra a exportação |
| GET | `/notificacoes` | Alertas de prazo vencido entre as tarefas visíveis |

`GET /tarefas`, `/tarefas/prioritaria`, `/metricas` e `/relatorios` aceitam filtros opcionais `setor_id`, `servidor_id`, `prioridade`, `status`, `etapa`, `tipo_conteudo`, `canal`, `busca`, `de` e `ate`. Datas em `AAAA-MM-DD`, intervalo inclusivo por data de criação UTC. O frontend usa filtros locais após carregar apenas as tarefas autorizadas. Relatórios PDF usam a impressão do navegador, sem endpoint de PDF.

Criação de tarefa:

```json
{
  "titulo": "Vídeo da campanha de vacinação",
  "tipo_conteudo": "video",
  "etapa": "pauta",
  "canal": "YouTube",
  "descricao": "Editar entrevistas, inserir legendas e revisar as informações oficiais.",
  "origem": "Protocolo",
  "prioridade": "alta",
  "setor_id": "uuid-do-setor",
  "servidor_id": null,
  "revisor_id": null,
  "link_material": "https://example.com/versao-1",
  "prazo_etapa": "2026-10-19T15:00:00Z",
  "publicacao_prevista": "2026-10-20T19:00:00Z",
  "prazo": "2026-10-20T18:00:00Z"
}
```

PUT envia o registro completo, a nova `etapa` e a `versao` recebida no GET. `anotacao` registra o motivo dos ajustes/reabertura e é obrigatória nas devoluções da aprovação e na reabertura de entregas publicadas. `status` é derivado da etapa; não pode ser usado para pular revisão ou aprovação. Status: `pendente`, `em_andamento`, `concluida`. Prioridades: `urgente`, `alta`, `media`, `baixa`. Origens: `Solicitação interna`, `Protocolo`, `Ouvidoria`.

Setor: `{nome, capacidade, cor}`. Usuário: `{nome, email, senha, setor_id, perfil}`; senha obrigatória na criação e opcional na edição. Perfil: `gestor` ou `servidor`. Configuração: `{pontos_por_tarefa: 10}` (1–1000).

O chat é consultado a cada 5 segundos enquanto o modal está aberto. Os avisos da interface mostram atrasos; não há envio de e-mail nem notificações push no MVP.

## Contrato editorial

Tipos: `video`, `materia`, `arte`, `outro`. Etapas: `pauta`, `producao`, `edicao`, `revisao`, `aprovacao`, `pronto_publicacao`, `publicado`; `encerrada` preserva demandas antigas concluídas. Regras detalhadas em [EDITORIAL.md](EDITORIAL.md).

A transição `aprovacao` → `pronto_publicacao` exige gestor do mesmo setor da demanda. A API gera `aprovador_id` e `aprovado_em`; valores enviados pelo cliente são ignorados. A transição para `publicado` exige aprovação válida e `link_publicado`, gerando `publicado_em`. Conteúdo alterado após aprovação retorna à revisão e limpa a assinatura. Concorrência: versão desatualizada retorna 409 sem alteração nem novo histórico.

Histórico inclui `etapa`, `prazo_etapa` e `responsavel_id` como fotografia da etapa no momento da operação. Métricas incluem `etapas`, `tipos_conteudo`, `publicadas` e `aguardando_aprovacao`. O CSV inclui canal, prazos e dados de aprovação/publicação.
