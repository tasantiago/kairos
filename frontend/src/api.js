import { setoresIniciais, usuariosIniciais, tarefasIniciais } from './data.js';
import { applyEditorial, etapaLabels, deadline } from './editorial.js';
const env = import.meta.env || {};
export const demo = env.VITE_DEMO_MODE !== 'false';
let token = '';
let demoUser = usuariosIniciais[0];
const demoFiles = new Map();
export function setDemoUser(user) { demoUser = user; }
const store = { tarefas: structuredClone(tarefasIniciais), setores: structuredClone(setoresIniciais), usuarios: structuredClone(usuariosIniciais), historicos: [], mensagens: [], config: { pontos_por_tarefa: 10 } };
let nextCode = Math.max(100, ...store.tarefas.map(t => Number(t.codigo.slice(4)))) + 1;
for (const t of store.tarefas) {
  store.historicos.push({ id: crypto.randomUUID(), tarefa_id: t.id, usuario_id: t.servidor_id, de_setor: null, para_setor: t.setor_id, etapa: t.etapa, prazo_etapa: deadline(t), responsavel_id: t.servidor_id, anotacao: `Demanda editorial · ${etapaLabels[t.etapa]}.`, data: t.criado_em });
  if (t.aprovado_em) store.historicos.push({ id: crypto.randomUUID(), tarefa_id: t.id, usuario_id: t.aprovador_id, para_setor: t.setor_id, anotacao: 'Versão revisada e aprovada para publicação.', data: t.aprovado_em });
  if (t.publicado_em) store.historicos.push({ id: crypto.randomUUID(), tarefa_id: t.id, usuario_id: t.servidor_id, para_setor: t.setor_id, anotacao: `Publicação registrada: ${t.link_publicado}`, data: t.publicado_em });
}
export async function api(path, options = {}) {
  if (!demo) {
    const response = await fetch((env.VITE_API_URL || '/api') + path, { ...options, headers: { ...(!(options.body instanceof FormData) ? { 'Content-Type': 'application/json' } : {}), ...(token ? { Authorization: `Bearer ${token}` } : {}), ...options.headers } });
    if (!response.ok) { const error = await response.json().catch(() => ({})); if (response.status === 401 && path !== '/login') { token = ''; window.dispatchEvent(new Event('fluxo:session-expired')); } const failure = new Error(error.erro || 'Não foi possível concluir a operação.'); failure.status = response.status; throw failure; }
    if (response.status === 204) return null;
    return response.json();
  }
  const body = typeof options.body === 'string' ? JSON.parse(options.body) : options.body;
  const method = options.method || 'GET';
  const [type, id, action] = path.replace(/^\//, '').split('/');
  const targetTask = type === 'tarefas' && id ? store.tarefas.find(t => t.id === id) : type === 'anexos' && id ? store.tarefas.find(t => t.id === demoFiles.get(id)?.tarefa_id) : null;
  if (targetTask && demoUser.perfil !== 'gestor' && targetTask.setor_id !== demoUser.setor_id && targetTask.servidor_id !== demoUser.id && targetTask.revisor_id !== demoUser.id) throw new Error('Você não tem acesso a esta demanda.');
  if (method !== 'GET' && ['setores', 'usuarios', 'config'].includes(type) && demoUser.perfil !== 'gestor') throw new Error('Esta ação é exclusiva de gestores.');
  if (method === 'DELETE' && ['tarefas', 'anexos'].includes(type) && demoUser.perfil !== 'gestor') throw new Error('Esta ação é exclusiva de gestores.');
  if (type === 'tarefas' && action === 'anexos') {
    if (method === 'POST') { const file = body.get('arquivo'); const fileId = crypto.randomUUID(); demoFiles.set(fileId, { file, tarefa_id: id, id: fileId, nome: file.name, tamanho: file.size }); }
    return [...demoFiles.values()].filter(f => f.tarefa_id === id).map(({ file, ...meta }) => meta);
  }
  if (type === 'anexos' && method === 'DELETE') { demoFiles.delete(id); return null; }
  if (type === 'config') { if (method === 'PUT') store.config = { ...store.config, ...body }; return structuredClone(store.config); }
  if (type === 'login') return { usuario: store.usuarios[0], token: 'demonstracao' };
  if (type === 'tarefas' && action === 'historico') return structuredClone(store.historicos.filter(h => h.tarefa_id === id));
  if (type === 'tarefas' && action === 'chat') {
    if (method === 'POST') store.mensagens.push({ id: crypto.randomUUID(), tarefa_id: id, usuario_id: demoUser.id, mensagem: body.mensagem, data: new Date().toISOString() });
    return structuredClone(store.mensagens.filter(m => m.tarefa_id === id));
  }
  if (!store[type]) throw new Error('Operação indisponível.');
  if (method === 'GET') return structuredClone(id ? store[type].find(item => item.id === id) : type === 'tarefas' && demoUser.perfil !== 'gestor' ? store.tarefas.filter(t => t.setor_id === demoUser.setor_id || t.servidor_id === demoUser.id || t.revisor_id === demoUser.id) : store[type]);
  if (type === 'tarefas' && body?.servidor_id && store.usuarios.find(u => u.id === body.servidor_id)?.setor_id !== body.setor_id) throw new Error('O responsável deve pertencer ao setor selecionado.');
  if (type === 'tarefas' && body?.revisor_id && store.usuarios.find(u => u.id === body.revisor_id)?.setor_id !== body.setor_id) throw new Error('O revisor deve pertencer ao setor selecionado.');
  if (type === 'usuarios' && id === demoUser.id && (method === 'DELETE' || body?.perfil !== 'gestor')) throw new Error('Você não pode excluir sua conta ou retirar seu próprio acesso de gestor.');
  if (method === 'POST') {
    if (type === 'tarefas' && demoUser.perfil !== 'gestor' && body.setor_id !== demoUser.setor_id) throw new Error('Crie demandas apenas no seu setor.');
    const validated = type === 'tarefas' ? applyEditorial(null, body, demoUser) : body;
    const created = { ...validated, id: crypto.randomUUID(), ...(type === 'tarefas' ? { codigo: `FLX-${String(nextCode++).padStart(4, '0')}`, criado_em: new Date().toISOString(), atualizado_em: new Date().toISOString() } : {}) }; delete created.senha;
    store[type].push(created);
    if (type === 'tarefas') store.historicos.push({ id: crypto.randomUUID(), tarefa_id: created.id, usuario_id: demoUser.id, para_setor: created.setor_id, etapa: created.etapa, prazo_etapa: deadline(created), responsavel_id: created.servidor_id, anotacao: 'Pauta criada.', data: created.criado_em });
    return structuredClone(created);
  }
  const index = store[type].findIndex(item => item.id === id);
  if (index < 0) throw new Error('Registro não encontrado.');
  if (method === 'DELETE') {
    if (type === 'setores' && (store.tarefas.some(t => t.setor_id === id) || store.usuarios.some(u => u.setor_id === id) || store.historicos.some(h => h.de_setor === id || h.para_setor === id))) throw new Error('Este setor está em uso.');
    if (type === 'usuarios' && (store.tarefas.some(t => t.servidor_id === id || t.revisor_id === id || t.aprovador_id === id) || store.historicos.some(h => h.usuario_id === id) || store.mensagens.some(m => m.usuario_id === id))) throw new Error('Esta pessoa possui tarefas ou histórico.');
    if (type === 'tarefas') { store.historicos = store.historicos.filter(h => h.tarefa_id !== id); store.mensagens = store.mensagens.filter(m => m.tarefa_id !== id); for (const [key, f] of demoFiles) if (f.tarefa_id === id) demoFiles.delete(key); }
    store[type].splice(index, 1); return null;
  }
  const previous = store[type][index];
  if (type === 'usuarios' && previous.setor_id !== body.setor_id && store.tarefas.some(t => t.servidor_id === id || t.revisor_id === id || t.aprovador_id === id)) throw new Error('Reatribua as tarefas antes de mudar o setor.');
  const validated = type === 'tarefas' ? applyEditorial(previous, { ...previous, ...body, versao: body.versao }, demoUser) : body;
  store[type][index] = { ...previous, ...validated, atualizado_em: new Date().toISOString() };
  delete store[type][index].senha;
  if (type === 'tarefas') { const current = store[type][index]; const movement = previous.etapa !== current.etapa ? `${etapaLabels[previous.etapa]} → ${etapaLabels[current.etapa]}.` : 'Dados editoriais atualizados.'; store.historicos.push({ id: crypto.randomUUID(), tarefa_id: id, usuario_id: demoUser.id, de_setor: previous.setor_id, para_setor: current.setor_id, etapa: current.etapa, prazo_etapa: deadline(current), responsavel_id: current.servidor_id, anotacao: `${movement}${body.anotacao ? ' ' + body.anotacao : ''}${previous.aprovado_em && !current.aprovado_em ? ' Aprovação invalidada; nova revisão necessária.' : ''}${previous.etapa !== 'publicado' && current.etapa === 'publicado' ? ' ' + current.link_publicado : ''}`, data: new Date().toISOString() }); }
  return structuredClone(store[type][index]);
}
export async function login(email, senha) { const result = await api('/login', { method: 'POST', body: JSON.stringify({ email, senha }) }); token = result.token; return result.usuario; }
export function logout() { token = ''; }
export async function downloadAttachment(id) {
  let blob;
  if (demo) { blob = demoFiles.get(id)?.file; if (!blob) throw new Error('Arquivo não encontrado.'); }
  else { const response = await fetch((env.VITE_API_URL || '/api') + '/anexos/' + id, { headers: { Authorization: `Bearer ${token}` } }); if (!response.ok) throw new Error('Não foi possível baixar o arquivo.'); blob = await response.blob(); }
  return blob;
}
