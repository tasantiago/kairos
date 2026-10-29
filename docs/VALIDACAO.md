# Validação

Neste ambiente, o acesso necessário para instalar dependências e publicar foi interrompido. O binário `go` disponível não é o compilador Go. Por isso não foi possível compilar a API, construir o React, executar os testes Go, conectar ao PostgreSQL ou testar visualmente a interface.

Checks locais executados com sucesso:

- Treze casos de teste Node, incluindo oito cenários do fluxo editorial: integridade entre responsáveis/setores, priorização, CSV seguro, ciclo completo de demanda/chat/arquivo, escopo/permissões, revisão por responsável, aprovação apenas por gestor do setor, edição obrigatória para vídeos, link de publicação obrigatório, invalidação da aprovação, reabertura justificada, concorrência e cronologia dos prazos.
- `node --check` dos módulos de dados, API e configuração Vite.
- Conferência de tags/aninhamento JSX com transformação local e check sintático do JavaScript resultante. Esse check não substitui o build Vite nem verifica as bibliotecas importadas.
- Conferência de arquivos JSON e referências relativas dos componentes.
- `git diff --cached --check` e confirmação de que os três `.env` estão ignorados pelo Git.

Os testes detectaram reutilização de código após exclusão de tarefa na demonstração. A sequência foi corrigida e todos os casos passaram após a correção.

Para validar antes de produção:

```bash
cd frontend
npm install
node --test src/data.test.js src/api.test.js src/editorial.test.js
npm run build
cd ../backend
go mod tidy
gofmt -w .
go vet ./...
go test ./...
cd ..
docker compose --env-file deploy/.env -f deploy/compose.yaml up --build -d
```

Roteiro integrado:

- Confirmar `/health` 200 e login com o administrador de bootstrap.
- Criar dois setores e servidores; confirmar CRUD e mensagens de conflito por nome/e-mail.
- Criar tarefa com responsável do setor errado: deve retornar 400 sem inserir tarefa ou histórico.
- Criar demanda válida; avançar por produção, edição, revisão, aprovação e publicação. Confirmar histórico com etapa, responsável, prazo e anotação.
- Confirmar que vídeos não pulam edição e que servidores/gestores de outros setores não aprovam.
- Publicar sem link ou aprovação deve falhar. Editar briefing/link após aprovação deve retornar à revisão e limpar assinatura.
- Gravar uma versão antiga deve retornar 409 sem mudar a demanda ou o histórico.
- Devolver da aprovação sem motivo deve falhar; com motivo, registrar a correção.
- Executar `002_editorial.sql` em banco com demandas antigas e confirmar que concluídas ficam em Encerrada (legado), sem falsa publicação.
- Entrar como servidor: listar apenas tarefas pessoais/do setor; negar leitura/edição de tarefa fora do escopo e administração de equipe/configurações.
- Reabrir publicação como gestor, com justificativa: limpar timestamps/assinatura e retirar os pontos correspondentes. Um servidor não pode reabrir entrega publicada.
- Enviar chat e arquivo; recarregar, entrar novamente e confirmar persistência. Arquivo maior que 10 MB deve falhar.
- Excluir tarefa como gestor e verificar remoção das mensagens, histórico, metadados e arquivo do volume.
- Exportar CSV com título começando por `=`: a célula precisa ser neutralizada para não executar fórmula.
- Verificar filtros, impressão/PDF, modais por teclado, telas de 375/768/1440 px e ampliação de texto de 200%.
- Reiniciar containers e confirmar dados preservados. Simular indisponibilidade do banco e verificar health 503.
- Confirmar JWT expirado/inválido 401, algoritmo diferente rejeitado, CORS para origem inesperada 403.

As migrations e bootstrap usam advisory locks. Para verificar concorrência, iniciar duas instâncias simultaneamente e confirmar uma aplicação de migration e um administrador.

Os testes Go de regra editorial foram acrescentados, mas não executados por ausência do compilador. Os treze casos Node validam a regra e o modo de demonstração; não validam PostgreSQL, build React ou interação real de navegador.
