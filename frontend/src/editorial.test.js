import test from 'node:test';
import assert from 'node:assert/strict';
import { applyEditorial, availableStages } from './editorial.js';
const manager = { id: 'gestor-ascom', perfil: 'gestor', setor_id: 'ascom' };
const editor = { id: 'editor', perfil: 'servidor', setor_id: 'ascom' };
const reviewer = { id: 'revisor', perfil: 'servidor', setor_id: 'ascom' };
const now = new Date('2026-10-09T14:00:00Z');
const brief = () => ({ titulo: 'Vídeo da campanha', descricao: 'Roteiro e legendas', tipo_conteudo: 'video', canal: 'YouTube', etapa: 'pauta', setor_id: 'ascom', revisor_id: reviewer.id, prazo: '2026-10-12T14:00:00Z', prazo_etapa: '2026-10-10T14:00:00Z', publicacao_prevista: '2026-10-12T17:00:00Z', link_material: 'https://example.com/version-1' });
const move = (current, etapa, actor = manager, extra = {}) => applyEditorial(current, { ...current, etapa, ...extra }, actor, now);
function approvalQueue(type = 'video') { let current = applyEditorial(null, { ...brief(), tipo_conteudo: type }, manager, now); for (const stage of type === 'video' ? ['producao', 'edicao', 'revisao'] : ['producao', 'revisao']) current = move(current, stage, editor); return move(current, 'aprovacao', reviewer); }

test('vídeo passa por edição, revisão, aprovação do gestor e publicação comprovada', () => {
  let current = approvalQueue(); current = move(current, 'pronto_publicacao');
  assert.equal(current.aprovador_id, manager.id); assert.equal(current.aprovado_em, now.toISOString());
  current = move(current, 'publicado', editor, { link_publicado: 'https://example.com/video-final' });
  assert.equal(current.status, 'concluida'); assert.ok(current.publicado_em); assert.ok(current.concluido_em);
});
test('matéria pode ir da redação à revisão sem exigir edição audiovisual', () => {
  const current = approvalQueue('materia'); assert.equal(current.etapa, 'aprovacao'); assert.equal(current.status, 'em_andamento');
});
test('servidor e gestor de outro setor não podem aprovar; campo forjado não assina', () => {
  const current = approvalQueue(); assert.throws(() => move(current, 'pronto_publicacao', editor), /gestor/); assert.throws(() => move(current, 'pronto_publicacao', { ...manager, setor_id: 'gabinete' }), /setor/);
  const forged = move(current, 'aprovacao', editor, { aprovado_em: now.toISOString(), aprovador_id: editor.id }); assert.equal(forged.aprovado_em, null);
  assert.ok(!availableStages(current, editor).includes('pronto_publicacao'));
});
test('não se pode pular edição de vídeo, revisão ou prova de publicação', () => {
  const current = move(applyEditorial(null, brief(), manager, now), 'producao');
  assert.throws(() => move(current, 'revisao'), /edição/); assert.throws(() => move(current, 'publicado'), /etapas/);
  const ready = move(approvalQueue(), 'pronto_publicacao'); assert.throws(() => move(ready, 'publicado', editor), /link/);
});
test('conteúdo alterado perde aprovação e retorna para revisão', () => {
  const ready = move(approvalQueue(), 'pronto_publicacao'); const modified = move(ready, ready.etapa, editor, { link_material: 'https://example.com/version-2' });
  assert.equal(modified.etapa, 'revisao'); assert.equal(modified.aprovado_em, null); assert.equal(modified.aprovador_id, null);
  assert.throws(() => move(ready, 'publicado', editor, { titulo: 'Novo conteúdo', link_publicado: 'https://example.com/publicacao' }), /nova revisão/);
  assert.equal(move(approvalQueue(), 'aprovacao', editor, { descricao: 'Texto atualizado' }).etapa, 'revisao');
});
test('revisão concluída pelo responsável e pedidos de ajuste documentados', () => {
  let current = applyEditorial(null, brief(), manager, now); for (const stage of ['producao', 'edicao', 'revisao']) current = move(current, stage, editor);
  assert.throws(() => move(current, 'aprovacao', editor), /revisor/);
  current = move(current, 'aprovacao', reviewer); assert.throws(() => move(current, 'edicao', manager), /motivo/);
  assert.equal(move(current, 'edicao', manager, { anotacao: 'Corrigir o áudio da entrevista.' }).etapa, 'edicao');
});
test('reabrir publicação exige gestor e justificativa; timestamps são limpos', () => {
  const published = move(move(approvalQueue(), 'pronto_publicacao'), 'publicado', editor, { link_publicado: 'https://example.com/video' });
  assert.throws(() => move(published, 'revisao', editor, { anotacao: 'Correção' }), /gestor/);
  assert.throws(() => move(published, 'revisao', manager), /motivo/);
  const revised = move(published, 'revisao', manager, { anotacao: 'Atualização de informação oficial.' }); assert.equal(revised.publicado_em, null); assert.equal(revised.concluido_em, null); assert.equal(revised.link_publicado, '');
});
test('edição concorrente, cronologia impossível e URLs inseguras são rejeitadas', () => {
  const current = approvalQueue(); assert.throws(() => move(current, 'aprovacao', manager, { versao: current.versao - 1 }), /mudou/);
  assert.throws(() => applyEditorial(null, { ...brief(), prazo_etapa: '2026-10-13T00:00:00Z' }, manager, now), /prazo da etapa/);
  assert.throws(() => applyEditorial(null, { ...brief(), publicacao_prevista: '2026-10-11T00:00:00Z' }, manager, now), /entrega/);
  assert.throws(() => applyEditorial(null, { ...brief(), link_material: 'javascript:alert(1)' }, manager, now), /HTTP/);
});
