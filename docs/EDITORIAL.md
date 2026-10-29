# Comunicação Social: fluxo editorial

O Fluxo é voltado à gestão da produção de vídeos, matérias e peças de Comunicação Social. Organiza o briefing, a edição, as revisões, a decisão do gestor e o prazo de publicação. A edição audiovisual é realizada pelas ferramentas da equipe; o sistema acompanha a etapa e a versão entregue para revisão.

## Etapas

| Etapa | Trabalho principal | Condição para avançar |
|---|---|---|
| Pauta | Definir objetivo, tipo, canal e briefing | Iniciar produção |
| Produção | Redação, roteiro, captação e informações oficiais | Vídeo vai à edição; matéria pode ir à revisão |
| Edição | Cortes, áudio, legendas, identidade visual ou preparação da peça | Entregar versão para revisão |
| Revisão | Conferir texto, dados, imagem, áudio e acessibilidade | Revisor designado ou gestor envia para aprovação |
| Aprovação | Decisão final dos gestores da Comunicação Social | Gestor do setor responsável aprova ou devolve com motivo |
| Pronto para publicar | Conteúdo aprovado aguardando execução no canal | Registrar o link da publicação realizada |
| Publicado | Entrega realizada e comprovada pelo link | Apenas gestor pode reabrir, justificando a correção |

O gestor pode pedir ajustes na edição ou revisão; não precisa aprovar uma versão com pendências. A transição para aprovação não pode ser feita por outro servidor quando há um revisor designado. Se ainda não houver revisor designado, a equipe do setor pode encaminhar a revisão; recomenda-se definir o revisor em cada demanda.

## Responsáveis e versões

- `servidor_id`: pessoa responsável pelo trabalho na etapa atual. Pode mudar ao passar da redação/edição para a publicação.
- `revisor_id`: pessoa que confere a entrega antes da aprovação.
- `aprovador_id` e `aprovado_em`: gerados pelo sistema somente quando um gestor do setor aprova.
- `link_material`: link da versão de trabalho, pasta compartilhada ou material para revisão. Para vídeos, prefira um link específico da versão entregue.
- `link_publicado`: URL da matéria, vídeo ou peça publicada no canal.
- `versao`: controla alterações concorrentes do registro. Uma gravação com versão desatualizada retorna 409 e exige reabrir a demanda.

Mudanças no título, briefing, tipo, canal, setor ou link do material invalidam uma aprovação existente e retornam a demanda à revisão. Se o conteúdo mudar enquanto aguarda aprovação, também volta à revisão. Alterar conteúdo e aprovar/publicar na mesma operação é rejeitado. Ajustes de prioridade, atribuição e prazo não mudam o conteúdo aprovado.

O sistema não monitora alterações feitas diretamente em arquivos ou links externos. Para uma nova edição, troque o link para a versão atual e registre a revisão. Anexos podem conter documentos e referências; não representam versionamento de vídeos ou aprovação automática dos arquivos.

## Prazos

| Campo | Significado |
|---|---|
| `prazo_etapa` | Data e hora para terminar a etapa atual; opcional, usa o prazo de entrega como referência se vazio |
| `prazo` | Prazo final de entrega da demanda; obrigatório |
| `publicacao_prevista` | Data e hora prevista para publicação no canal; opcional |

Quando informados, a ordem precisa ser **prazo da etapa ≤ prazo de entrega ≤ publicação prevista**. Datas passadas são permitidas para representar atrasos reais. Ao avançar de etapa, o usuário atualiza responsável e prazo conforme o planejamento da equipe.

Cada alteração salva no histórico a etapa, o responsável e o prazo vigente. A prioridade considera urgência, origem, tempo sem movimentação, atraso da etapa/publicação e proximidade de até 24 horas do prazo da etapa (ou da entrega se não houver prazo de etapa).

O painel mostra a fila por etapa e as próximas publicações por data e canal. O horário escolhido é interpretado no fuso do navegador e armazenado como timestamp; relatórios da API exportam em RFC3339.

## Publicação e preservação de dados antigos

O sistema registra publicações realizadas pela equipe; não envia conteúdos automaticamente ao portal, YouTube ou redes sociais. “Pronto para publicar” não significa que o conteúdo já esteja no ar.

A migration `002_editorial.sql` acrescenta campos e regras sem reescrever a migration inicial. Demandas antigas concluídas ficam em **Encerrada (legado)** para não serem tratadas como publicações aprovadas. Registros antigos em andamento viram Produção; os pendentes, Pauta. O bootstrap de bancos novos cria Comunicação Social; bancos existentes preservam seus setores e usuários.
