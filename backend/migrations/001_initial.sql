CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE setores (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), nome varchar(120) UNIQUE NOT NULL,
 capacidade integer NOT NULL DEFAULT 10 CHECK (capacidade BETWEEN 1 AND 10000),
 cor varchar(7) NOT NULL DEFAULT '#32886c'
);
CREATE TABLE usuarios (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), nome varchar(120) NOT NULL,
 email varchar(254) UNIQUE NOT NULL, senha_hash text NOT NULL,
 setor_id uuid NOT NULL REFERENCES setores(id) ON DELETE RESTRICT,
 perfil varchar(10) NOT NULL DEFAULT 'servidor' CHECK (perfil IN ('gestor','servidor')),
 criado_em timestamptz NOT NULL DEFAULT now()
);
CREATE SEQUENCE tarefas_codigo_seq START 101;
CREATE FUNCTION fluxo_codigo() RETURNS text LANGUAGE sql VOLATILE AS $$
 SELECT 'FLX-' || lpad(n::text,GREATEST(4,length(n::text)),'0') FROM (SELECT nextval('tarefas_codigo_seq') AS n) s
$$;
CREATE TABLE tarefas (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 codigo text UNIQUE NOT NULL DEFAULT fluxo_codigo(),
 titulo varchar(200) NOT NULL, descricao text NOT NULL DEFAULT '', origem varchar(40) NOT NULL,
 status varchar(20) NOT NULL DEFAULT 'pendente' CHECK (status IN ('pendente','em_andamento','concluida')),
 prioridade varchar(10) NOT NULL CHECK (prioridade IN ('urgente','alta','media','baixa')),
 servidor_id uuid REFERENCES usuarios(id) ON DELETE RESTRICT,
 setor_id uuid NOT NULL REFERENCES setores(id) ON DELETE RESTRICT,
 prazo timestamptz NOT NULL, criado_em timestamptz NOT NULL DEFAULT now(),
 atualizado_em timestamptz NOT NULL DEFAULT now(), concluido_em timestamptz,
 CHECK ((status = 'concluida') = (concluido_em IS NOT NULL))
);
CREATE INDEX idx_tarefas_setor_status ON tarefas(setor_id,status);
CREATE INDEX idx_tarefas_servidor ON tarefas(servidor_id);
CREATE INDEX idx_tarefas_prazo ON tarefas(prazo) WHERE status <> 'concluida';
CREATE TABLE historicos_tarefa (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), tarefa_id uuid NOT NULL REFERENCES tarefas(id) ON DELETE CASCADE,
 de_setor uuid REFERENCES setores(id) ON DELETE RESTRICT, para_setor uuid REFERENCES setores(id) ON DELETE RESTRICT,
 usuario_id uuid REFERENCES usuarios(id) ON DELETE RESTRICT, anotacao text NOT NULL, data timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_historicos_tarefa_data ON historicos_tarefa(tarefa_id,data);
CREATE TABLE chat_tarefa (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), tarefa_id uuid NOT NULL REFERENCES tarefas(id) ON DELETE CASCADE,
 usuario_id uuid REFERENCES usuarios(id) ON DELETE RESTRICT, mensagem text NOT NULL CHECK (length(mensagem) BETWEEN 1 AND 4000),
 data timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_chat_tarefa_data ON chat_tarefa(tarefa_id,data);
CREATE TABLE anexos (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), tarefa_id uuid NOT NULL REFERENCES tarefas(id) ON DELETE CASCADE,
 usuario_id uuid REFERENCES usuarios(id) ON DELETE RESTRICT, nome varchar(255) NOT NULL,
 mime text NOT NULL, tamanho bigint NOT NULL, arquivo text NOT NULL, criado_em timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_anexos_tarefa ON anexos(tarefa_id);
CREATE TABLE relatorios (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), usuario_id uuid REFERENCES usuarios(id) ON DELETE RESTRICT,
 filtros jsonb NOT NULL DEFAULT '{}', formato text NOT NULL CHECK (formato IN ('csv','pdf')),
 criado_em timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE configuracoes (id integer PRIMARY KEY CHECK (id=1), pontos_por_tarefa integer NOT NULL DEFAULT 10 CHECK (pontos_por_tarefa BETWEEN 1 AND 1000));
