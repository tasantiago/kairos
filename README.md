# Fluxo · Comunicação Social

Gestão da produção editorial da Comunicação Social: vídeos, matérias e peças gráficas, da pauta à revisão, aprovação pelos gestores da área e publicação. O MVP reúne briefing, responsáveis, canal e prazos em português.

## Organização neste repositório

Esta entrega foi importada para `tasantiago/kairos`. As pastas anteriores `api/` e `front/`, assim como o `docker-compose.yml` da raiz, foram preservados. A versão editorial descrita aqui usa `backend/`, `frontend/` e `deploy/compose.yaml`; siga os comandos abaixo para executar este MVP.

## O que foi implementado

- React + Vite, Tailwind v4, componentes no padrão Shadcn com Radix UI, React Router, Lucide e Recharts.
- Visão geral, tarefas em lista e quadro, tarefas pessoais, setores, equipe, relatórios, ranking e configuração de pontos.
- Quadro editorial: Pauta → Produção → Edição → Revisão → Aprovação → Pronto para publicar → Publicado. Vídeos passam pela edição; matérias podem seguir da redação à revisão.
- Briefing/roteiro, tipo de conteúdo, canal, responsável, revisor, link do material e link da publicação.
- Aprovação final pelos gestores do setor responsável. Conteúdo alterado perde a aprovação; publicação exige aprovação válida e link.
- Prazos com data e hora para a etapa, entrega e publicação prevista; fila editorial, próximas publicações e histórico com responsável e prazo de cada etapa.
- Controle de versões para impedir que uma alteração concorrente sobrescreva decisões editoriais.
- Chat com polling de 5 segundos e anexos de até 10 MB.
- Busca, filtros, CSV compatível com Excel e PDF através do diálogo de impressão do navegador.
- API Go + Gorilla Mux + pgx, autenticação JWT, senhas bcrypt, permissões por perfil e setor, CORS e logs estruturados.
- PostgreSQL, migrations transacionais, volume para anexos, Docker Compose e health checks.

## Estado da entrega

O código está preparado no workspace. A demonstração do frontend é explicitamente temporária: usa dados ilustrativos em memória e reinicia ao recarregar. Nenhum dado de trabalho é salvo no navegador. O modo integrado usa a API e persiste tarefas, usuários, setores, histórico, mensagens e configurações no PostgreSQL; anexos ficam no volume de uploads.

O acesso de rede para baixar dependências/publicar foi interrompido. Este ambiente não contém o compilador Go (o executável `go` disponível não é o compilador da linguagem) nem as dependências do frontend. **Build, testes Go, integração com PostgreSQL e QA de navegador não foram executados. Não há URL publicada.** Veja [docs/VALIDACAO.md](docs/VALIDACAO.md) para os checks realizados e o roteiro de validação restante. `go.sum` e `package-lock.json` devem ser gerados ao instalar as dependências e versionados após validação.

## Desenvolvimento

Requisitos: Go 1.23+, Node.js 22+, PostgreSQL 16+; Docker Compose é a alternativa mais simples.

Frontend em demonstração:

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Abra `http://localhost:5173`. O seletor de perfil permite explorar os painéis de gestor e servidor na demonstração.

Para executar a API:

```bash
cd backend
cp .env.example .env
# Preencha DATABASE_URL, JWT_SECRET, ADMIN_EMAIL e ADMIN_PASSWORD.
go mod tidy
go test ./...
go run .
```

O bootstrap exige uma senha de 8–72 bytes e um segredo JWT de pelo menos 32 caracteres. Geração de segredo: `openssl rand -hex 32`. Não use os valores descritivos do exemplo. O administrador é criado somente se o banco não contiver usuários. As migrations são aplicadas automaticamente, registradas em `schema_migrations` e protegidas por advisory lock. `AUTO_MIGRATE=false` desativa essa etapa.

Para conectar o frontend, altere `frontend/.env`:

```dotenv
VITE_API_URL=/api
VITE_DEMO_MODE=false
```

Reinicie o Vite. O proxy encaminha `/api` para a API em `localhost:8080`. Faça login com as credenciais de bootstrap. O JWT fica apenas em memória; recarregar ou sair exige novo login.

## Execução com Docker

```bash
cp deploy/.env.example deploy/.env
# Preencha os campos vazios. Para POSTGRES_PASSWORD, use um segredo hexadecimal
# (ou faça percent-encoding se houver caracteres especiais na URI de conexão).
docker compose --env-file deploy/.env -f deploy/compose.yaml up --build -d
```

Abra `http://localhost:8088`. Essa execução sempre integra o frontend à API real. O PostgreSQL não expõe porta pública. Os volumes preservam banco e anexos. `docker compose down` preserva os volumes; `down -v` apaga os dados.

Em ambientes com proxy e CA de sessão, forneça a CA de confiança aos builds com `--secret id=proxy_ca,src=<caminho-da-ca>`, conforme a política da infraestrutura. Os Dockerfiles aceitam o secret opcional e mantêm a verificação TLS habilitada. O deploy em VPS precisa de domínio/HTTPS e configuração do `CORS_ORIGIN` para a origem real.

## Documentação

- [Fluxo editorial da Comunicação Social](docs/EDITORIAL.md)
- [Arquitetura e regras](docs/ARQUITETURA.md)
- [API](docs/API.md)
- [Manual do usuário](docs/MANUAL.md)
- [Roadmap das seis fases](docs/ROADMAP.md)
- [Validação e limitações](docs/VALIDACAO.md)

## Estrutura

```text
backend/     main.go, handlers, services, repositories, models, utils, migrations
frontend/    React, componentes UI, dados de demonstração, estilos
docs/        escopo, arquitetura, API, manual, validação e roadmap
deploy/      Dockerfiles, nginx, Compose e configuração de exemplo
```
