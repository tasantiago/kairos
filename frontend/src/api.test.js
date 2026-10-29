import test from 'node:test';
import assert from 'node:assert/strict';
import { api, setDemoUser, downloadAttachment } from './api.js';
import { usuariosIniciais, setoresIniciais } from './data.js';
const write = (method, body) => ({ method, body: JSON.stringify(body) });

test('demanda editorial, revisão e aprovação mantêm histórico, conversa e arquivo', async () => {
  const manager = usuariosIniciais[0]; setDemoUser(manager);
  const input = { titulo: 'Matéria do teste', descricao: 'Contexto', tipo_conteudo: 'materia', canal: 'Portal institucional', origem: 'Protocolo', prioridade: 'alta', setor_id: setoresIniciais[0].id, servidor_id: usuariosIniciais[1].id, revisor_id: usuariosIniciais[4].id, prazo: new Date(Date.now() + 3 * 86400000).toISOString() };
  const externalSector = await api('/setores', write('POST', { nome: 'Outro setor', capacidade: 5, cor: '#607ace' }));
  const externalUser = await api('/usuarios', write('POST', { nome: 'Outro responsável', email: 'outro@example.com', perfil: 'servidor', setor_id: externalSector.id, senha: 'senha-temporaria' }));
  await assert.rejects(api('/tarefas', write('POST', { ...input, servidor_id: externalUser.id })), /responsável/);
  await api('/usuarios/' + externalUser.id, { method: 'DELETE' }); await api('/setores/' + externalSector.id, { method: 'DELETE' });
  let current = await api('/tarefas', write('POST', input)); const id = current.id; const codigo = current.codigo;
  assert.equal((await api(`/tarefas/${id}/historico`)).length, 1);
  const statusOnly = await api('/tarefas/' + id, write('PUT', { ...current, status: 'concluida' })); assert.equal(statusOnly.status, 'pendente'); current = statusOnly;
  for (const stage of ['producao', 'revisao', 'aprovacao', 'pronto_publicacao']) current = await api('/tarefas/' + id, write('PUT', { ...current, etapa: stage, anotacao: 'Entrega revisada' }));
  setDemoUser(usuariosIniciais[1]); current = await api('/tarefas/' + id, write('PUT', { ...current, etapa: 'publicado', link_publicado: 'https://example.com/materia' }));
  assert.ok(current.concluido_em); assert.ok((await api(`/tarefas/${id}/historico`)).some(h => h.anotacao.includes('Entrega revisada')));
  setDemoUser(manager); current = await api('/tarefas/' + id, write('PUT', { ...current, etapa: 'revisao', anotacao: 'Informação atualizada' })); assert.equal(current.concluido_em, null); assert.equal(current.aprovado_em, null);
  await assert.rejects(api('/tarefas/' + id, write('PUT', { ...current, versao: 1 })), /mudou/);
  await api(`/tarefas/${id}/chat`, write('POST', { mensagem: 'Precisamos confirmar a entrega.' })); assert.equal((await api(`/tarefas/${id}/chat`)).length, 1);
  const form = new FormData(); form.append('arquivo', new Blob(['documento de teste']), 'teste.txt'); const files = await api(`/tarefas/${id}/anexos`, { method: 'POST', body: form }); assert.equal(await (await downloadAttachment(files[0].id)).text(), 'documento de teste');
  await api('/tarefas/' + id, { method: 'DELETE' }); assert.equal((await api(`/tarefas/${id}/historico`)).length, 0); assert.equal((await api(`/tarefas/${id}/chat`)).length, 0); assert.equal((await api(`/tarefas/${id}/anexos`)).length, 0);
  const next = await api('/tarefas', write('POST', input)); assert.notEqual(codigo, next.codigo); await api('/tarefas/' + next.id, { method: 'DELETE' });
});
test('escopo, gestão e referências também são respeitados na demonstração', async () => {
  setDemoUser(usuariosIniciais[1]); const tasks = await api('/tarefas'); assert.ok(tasks.length); assert.ok(tasks.every(t => t.setor_id === usuariosIniciais[1].setor_id || t.servidor_id === usuariosIniciais[1].id || t.revisor_id === usuariosIniciais[1].id));
  const waiting = tasks.find(t => t.etapa === 'aprovacao'); await assert.rejects(api('/tarefas/' + waiting.id, write('PUT', { ...waiting, etapa: 'pronto_publicacao' })), /gestor/); await assert.rejects(api('/config', write('PUT', { pontos_por_tarefa: 500 })), /gestores/);
  setDemoUser(usuariosIniciais[0]); await assert.rejects(api('/setores/' + setoresIniciais[0].id, { method: 'DELETE' }), /em uso/); await assert.rejects(api('/usuarios/' + usuariosIniciais[0].id, { method: 'DELETE' }), /próprio/);
  const u = await api('/usuarios', write('POST', { nome: 'Pessoa de teste', email: 'teste@exemplo.local', perfil: 'servidor', setor_id: setoresIniciais[0].id, senha: 'segredo-temporario' })); assert.equal(u.senha, undefined); assert.equal((await api('/usuarios/' + u.id)).senha, undefined); await api('/usuarios/' + u.id, { method: 'DELETE' });
});
