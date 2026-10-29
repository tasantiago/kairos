import test from 'node:test';
import assert from 'node:assert/strict';
import { tarefasIniciais, setoresIniciais, usuariosIniciais, score, csv } from './data.js';

test('cada responsável da demonstração pertence ao setor da tarefa', () => {
  for (const t of tarefasIniciais) { assert.equal(usuariosIniciais.find(u => u.id === t.servidor_id).setor_id, t.setor_id); assert.ok(setoresIniciais.some(s => s.id === t.setor_id)); }
});
test('priorização respeita urgência, bônus limitado de espera e prazo', () => {
  const now = new Date();
  const base = { prioridade: 'media', origem: 'Protocolo', atualizado_em: now.toISOString(), prazo: new Date(+now + 3 * 86400000).toISOString(), status: 'pendente' };
  assert.ok(score({ ...base, prioridade: 'urgente' }) > score(base));
  assert.ok(Math.abs(score({ ...base, atualizado_em: new Date(+now - 999 * 86400000).toISOString() }) - score(base) - 30) < .001);
  assert.ok(Math.abs(score({ ...base, prazo: new Date(+now - 86400000).toISOString() }) - score(base) - 25) < .001);
  assert.ok(Math.abs(score({ ...base, status: 'concluida', prazo: new Date(+now - 86400000).toISOString() }) - score(base)) < .001);
});
test('CSV escapa aspas, preserva Unicode e neutraliza fórmulas com espaços', () => {
  const task = { ...tarefasIniciais[0], titulo: ' =HYPERLINK("teste")' };
  const exported = csv([task], setoresIniciais, usuariosIniciais);
  assert.ok(exported.startsWith('\uFEFF')); assert.ok(exported.includes('"\' =HYPERLINK(""teste"")"')); assert.ok(exported.includes('Comunicação Social'));
});
