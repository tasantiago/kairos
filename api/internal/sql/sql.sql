CREATE DATABASE sistematarefas;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS usuarios

CREATE TABLE usuarios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    senha TEXT NOT NULL,
    setor_id UUID NOT NULL,
    tipo TEXT NOT NULL CHECK (tipo IN ('servidor', 'gestor', 'admin')),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE setores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome TEXT NOT NULL UNIQUE,
    capacidade_ideal INTEGER NOT NULL,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tarefas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    titulo TEXT NOT NULL,
    descricao TEXT,
    origem TEXT NOT NULL CHECK (origem IN ('telefone', 'whatsapp', 'presencial', 'email', 'outro')),
    prioridade TEXT CHECK (prioridade IN ('baixa', 'media', 'alta', 'urgente')),
    status TEXT CHECK (status IN ('pendente', 'em_andamento', 'concluida')),
    setor_atual_id UUID NOT NULL,
    responsavel_id UUID,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE movimentacoes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tarefa_id UUID NOT NULL,
    de_setor_id UUID,
    para_setor_id UUID NOT NULL,
    responsavel_id UUID NOT NULL,
    data_movimentacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notificacoes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tarefa_id UUID,
    tipo TEXT NOT NULL CHECK (tipo IN ('atribuicao', 'urgencia', 'nova_tarefa', 'alerta')),
    conteudo TEXT NOT NULL,
    destino TEXT NOT NULL,
    lida BOOLEAN DEFAULT FALSE,
    criada_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE usuarios
ADD CONSTRAINT fk_usuario_setor
FOREIGN KEY (setor_id) REFERENCES setores(id);

ALTER TABLE tarefas
ADD CONSTRAINT fk_tarefa_setor
FOREIGN KEY (setor_atual_id) REFERENCES setores(id);

ALTER TABLE tarefas
ADD CONSTRAINT fk_tarefa_usuario
FOREIGN KEY (responsavel_id) REFERENCES usuarios(id);

ALTER TABLE movimentacoes
ADD FOREIGN KEY (tarefa_id) REFERENCES tarefas(id);

ALTER SYSTEM SET timezone TO 'America/Porto_Velho';
