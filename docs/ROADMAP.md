# Roadmap e limites do MVP

O roteiro original foi organizado como uma entrega inicial para Comunicação Social, com itens de produção e extensões separados. “Implementado” significa código preparado; execução integrada ainda depende da validação descrita em VALIDACAO.md.

| Fase | Implementado no código | Ainda pendente |
|---|---|---|
| 1 – Setup | Nome, objetivo, MVP, Git local, backend/frontend/docs/deploy, exemplos e arquivos locais de ambiente | Instalar dependências e compilador; gerar lockfiles; conectar repositório remoto |
| 2 – Backend | Gorilla Mux, JWT, CORS, logs, pgx, migrations, CRUD, histórico, filtros, métricas, recomendação e transações de revisão/aprovação/publicação | Executar testes e validar transações/autorização com PostgreSQL |
| 3 – Frontend | React/Vite/Tailwind, componentes Shadcn/Radix, rotas, login, painéis, lista/quadro editorial, briefing, prazos por etapa, fila de aprovação, próximas publicações, chat, timeline e mapa de setores | Build, QA real em navegador e revisão com usuários |
| 4 – Inovação | Pontos/ranking configurável, chat persistente, histórico, avisos visuais de atraso | Heatmap por horário e movimentações históricas; notificações de transferência; e-mails; agenda de alertas |
| 5 – Relatórios | Carga, média de criação até conclusão no setor atual, recorrência exata de títulos na API, CSV e PDF via navegador | Tempo de permanência por setor, agrupamento semântico de tarefas recorrentes, atualização automática do dashboard |
| 6 – Produção | Dockerfiles, Compose, volumes, health checks, logs, documentação e manual | Deploy em VPS, HTTPS, backups/restauração, observabilidade externa, alertas de uptime, treinamento real |

O chat atualiza em polling. Os indicadores atualizam ao entrar no workspace ou após operações; não há assinatura em tempo real. O mapa é agregado por criação, dia da semana e setor atual. As médias não atribuem intervalos aos setores pelos quais a tarefa passou. Não há integração com sistemas administrativos externos.

O upload tem limite de tamanho e entrega arquivos como download, mas não possui antivírus. A primeira produção deverá decidir armazenamento externo/backup, política de retenção e permissões detalhadas com base no fluxo real da instituição.

A edição audiovisual e a publicação nos canais continuam sendo executadas pelas ferramentas da equipe. O Fluxo coordena essas etapas e registra suas entregas. Integração com CMS/redes sociais, armazenamento e versionamento de vídeos grandes, pré-visualização de vídeo e notificações de pedidos de aprovação poderão ser acrescentados conforme a rotina real.
