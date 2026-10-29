# Arquitetura e escopo

**Nome:** Fluxo. **Objetivo:** coordenar a produção de Comunicação Social, com responsáveis, revisores, aprovação e prazos de publicação.

O MVP cobre pauta → produção → edição → revisão → aprovação → pronto para publicar → publicado. As regras de conteúdo, assinatura do gestor e prazos estão em [EDITORIAL.md](EDITORIAL.md). O gestor administra equipes, setores e pontuação. O servidor acessa tarefas do próprio setor ou atribuídas a ele e pode encaminhá-las a outro setor. Após encaminhamento, perde acesso se não pertencer mais ao setor nem for o responsável.

O frontend usa React Router para navegação, Recharts para indicadores, Radix para modais acessíveis e Tailwind com um tema próprio. Há estados de carregamento, erro, vazio e sucesso, comportamento responsivo, filtros e controles acessíveis por teclado.

A API é dividida em handlers (HTTP, validação e autorização), repositories (acesso ao banco e transações), services (regra de priorização), models (contratos JSON) e utils (respostas e erros). JWT usa HS256, issuer e audience fixos, expiração de 8 horas e algoritmo explicitamente permitido. O usuário e seu perfil são recarregados do banco a cada requisição. O login limita tentativas por endereço remoto a 20 por 15 minutos por instância; atrás de proxy, configure rate limiting também na borda para evitar que todos os usuários compartilhem o limite do IP do proxy.

O responsável e o revisor devem pertencer ao setor escolhido. Um gestor precisa pertencer ao setor responsável para aprovar seu conteúdo. Aprovação, histórico e dados da demanda são gravados na mesma transação, com bloqueio da linha e comparação de `versao`. Novos clientes devem enviar o registro completo e a versão recebida no último GET. Tarefa e histórico são gravados em uma única transação. O campo `concluido_em` permanece estável em edições de tarefas concluídas e é limpo ao reabrir. Usuários e setores referenciados são protegidos contra exclusão. O gestor não pode excluir a si mesmo nem remover seu próprio perfil de gestor. Mensagens e histórico pertencem à tarefa e são excluídos com ela.

Anexos são salvos com nomes UUID e permissões restritas no volume. O nome original serve apenas para identificação e download. Downloads exigem autorização na tarefa e são enviados como attachment/octet-stream com `nosniff`. Limite: 10 MB. Exclusões removem metadados e tentam remover o arquivo; falhas de limpeza são registradas em log e podem exigir manutenção do volume.

Prioridade recomendada: urgência (100/65/35/10) + origem (Ouvidoria 15, demais 5) + tempo parado (3/dia, limite 30) + atraso (25). Tarefas concluídas não entram na recomendação. A priorização também considera atraso da etapa/publicação e adiciona 10 pontos para o prazo da etapa/entrega nas próximas 24 horas. Os pontos do ranking são calculados por entrega finalizada, evitando duplicidade ao reabrir e concluir novamente.

O tempo médio é o intervalo entre criação e conclusão agrupado pelo setor atual da tarefa. Não mede permanência em cada setor anterior. O mapa na página Setores mostra tarefas criadas por dia da semana e setor; o painel principal destaca a fila editorial e próximas publicações; um mapa por horas e movimentações históricas permanece no roadmap.

Produção própria: navegador → Nginx → API Go → PostgreSQL; arquivos em volume separado. Banco e backend não expõem portas no Compose. HTTPS deve ser terminado pelo proxy da VPS. A prévia Sites não foi criada: Go e sockets PostgreSQL requerem infraestrutura compatível, e o acesso necessário à publicação foi interrompido.
