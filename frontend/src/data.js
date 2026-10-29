import { tipoLabels, etapaLabels, deadline, publicationLate, statusForStage } from './editorial.js';
export const setoresIniciais = [
  { id: '11111111-1111-4111-8111-111111111111', nome: 'Comunicação Social', capacidade: 20, cor: '#32886c' },
];
export const usuariosIniciais = [
  { id: 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa', nome: 'Mariana Costa', email: 'mariana@fluxo.local', perfil: 'gestor', setor_id: setoresIniciais[0].id },
  { id: 'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb', nome: 'Rafael Almeida', email: 'rafael@fluxo.local', perfil: 'servidor', setor_id: setoresIniciais[0].id },
  { id: 'cccccccc-cccc-4ccc-8ccc-cccccccccccc', nome: 'Camila Santos', email: 'camila@fluxo.local', perfil: 'servidor', setor_id: setoresIniciais[0].id },
  { id: 'dddddddd-dddd-4ddd-8ddd-dddddddddddd', nome: 'Pedro Oliveira', email: 'pedro@fluxo.local', perfil: 'servidor', setor_id: setoresIniciais[0].id },
  { id: 'eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee', nome: 'Ana Rodrigues', email: 'ana@fluxo.local', perfil: 'servidor', setor_id: setoresIniciais[0].id },
];
function relativeDate(days, hour = 10) { const d = new Date(); d.setDate(d.getDate() + days); d.setHours(hour, 0, 0, 0); return d.toISOString(); }
const briefs = [
  ['Vídeo de prestação de contas: cortes e legendas', 'Editar entrevistas, inserir legendas e finalizar a versão para revisão. Material bruto na pasta compartilhada da equipe.', 'video', 'edicao', 1, 'urgente', -1, 1, 'YouTube'],
  ['Matéria sobre a entrega da nova unidade de saúde', 'Revisar nomes, dados e declarações; conferir fotos e crédito antes de enviar à coordenação.', 'materia', 'revisao', 2, 'alta', 0, 1, 'Portal institucional'],
  ['Campanha de vacinação: vídeo de orientação', 'Aguardar aprovação da coordenação para liberar o corte final com Libras e legendas.', 'video', 'aprovacao', 3, 'alta', 0, 2, 'Instagram'],
  ['Publicação da matéria sobre inscrições abertas', 'Matéria revisada e aprovada. Conferir chamada e registrar o link após a publicação.', 'materia', 'pronto_publicacao', 2, 'alta', 1, 1, 'Portal institucional'],
  ['Cobertura audiovisual do encontro de servidores', 'Preparar roteiro e captar depoimentos, imagens de apoio e áudio ambiente.', 'video', 'producao', 1, 'media', 2, 4, 'YouTube'],
  ['Pauta: balanço das ações de assistência social', 'Definir fontes, abordagem e fotos necessárias para a matéria.', 'materia', 'pauta', 2, 'media', 3, 6, 'Portal institucional'],
  ['Card com os horários de atendimento ao público', 'Produzir peça para redes sociais e conferir endereços e horários.', 'arte', 'edicao', 3, 'urgente', -1, 0, 'Instagram'],
  ['Revisão do vídeo sobre coleta seletiva', 'Conferir áudio, acessibilidade, identidade visual e os dados do serviço.', 'video', 'revisao', 1, 'alta', 1, 3, 'YouTube'],
  ['Aprovação da matéria sobre obras no centro', 'Conferir a versão revisada com as informações oficiais do cronograma.', 'materia', 'aprovacao', 2, 'alta', 0, 2, 'Portal institucional'],
  ['Cobertura fotográfica da feira de cultura', 'Selecionar fotos e publicar a cobertura já aprovada.', 'materia', 'publicado', 4, 'media', -2, -1, 'Portal institucional'],
  ['Vídeo: serviços digitais para o cidadão', 'Entrega final publicada com legendas e identificação institucional.', 'video', 'publicado', 1, 'media', -3, -2, 'YouTube'],
  ['Arte: calendário de eventos do mês', 'Peça aprovada e publicada nos canais institucionais.', 'arte', 'publicado', 3, 'baixa', -2, -1, 'Instagram'],
];
export const tarefasIniciais = briefs.map((r, i) => ({
  id: `10000000-0000-4000-8000-${String(i + 1).padStart(12, '0')}`, codigo: `FLX-${String(i + 101).padStart(4, '0')}`,
  titulo: r[0], descricao: r[1], setor_id: setoresIniciais[0].id, servidor_id: usuariosIniciais[r[4]].id,
  tipo_conteudo: r[2], etapa: r[3], canal: r[8], prioridade: r[5], status: statusForStage(r[3]),
  prazo_etapa: relativeDate(r[6], 12), prazo: relativeDate(r[7], 14), publicacao_prevista: relativeDate(r[7], 17),
  origem: 'Solicitação interna', revisor_id: usuariosIniciais[4].id, link_material: '', versao: 1,
  link_publicado: r[3] === 'publicado' ? `https://example.com/publicacoes/fluxo-${i + 101}` : '',
  aprovador_id: ['pronto_publicacao', 'publicado'].includes(r[3]) ? usuariosIniciais[0].id : null,
  aprovado_em: ['pronto_publicacao', 'publicado'].includes(r[3]) ? relativeDate(-4) : null,
  publicado_em: r[3] === 'publicado' ? relativeDate(r[7]) : null,
  criado_em: relativeDate(-10 + i % 4, 8 + i % 8), atualizado_em: relativeDate(-i % 3),
  concluido_em: r[3] === 'publicado' ? relativeDate(r[7]) : null,
}));
export const statusLabels = { pendente: 'Pendente', em_andamento: 'Em andamento', concluida: 'Concluída' };
export const prioridadeLabels = { urgente: 'Urgente', alta: 'Alta', media: 'Média', baixa: 'Baixa' };
export const initials = (name = '') => name.split(' ').map(x => x[0]).slice(0, 2).join('');
export const isLate = t => t.status !== 'concluida' && (new Date(deadline(t)) < new Date() || publicationLate(t));
export function score(t) {
  const due = new Date(deadline(t));
  const dueSoon = due >= new Date() && due - Date.now() <= 86400000;
  return ({ urgente: 100, alta: 65, media: 35, baixa: 10 }[t.prioridade] || 0) + (t.origem === 'Ouvidoria' ? 15 : 5) + Math.min(30, Math.max(0, (Date.now() - new Date(t.atualizado_em)) / 86400000) * 3) + (isLate(t) ? 25 : 0) + (t.status !== 'concluida' && dueSoon ? 10 : 0);
}
export function csv(tasks, setores, usuarios) {
  const cell = v => { let text = String(v ?? ''); if (/^\s*[=+@\-]/.test(text) || /^[\t\r\n]/.test(text)) text = "'" + text; return '"' + text.replaceAll('"', '""') + '"'; };
  const date = value => value ? new Date(value).toLocaleString('pt-BR') : '';
  return '\uFEFF' + [['Código', 'Título', 'Tipo', 'Etapa', 'Canal', 'Setor', 'Responsável', 'Revisor', 'Prioridade', 'Prazo da etapa', 'Prazo de entrega', 'Publicação prevista', 'Aprovado por', 'Aprovado em', 'Link publicado'], ...tasks.map(t => [t.codigo, t.titulo, tipoLabels[t.tipo_conteudo], etapaLabels[t.etapa], t.canal, setores.find(s => s.id === t.setor_id)?.nome, usuarios.find(u => u.id === t.servidor_id)?.nome, usuarios.find(u => u.id === t.revisor_id)?.nome, prioridadeLabels[t.prioridade], date(t.prazo_etapa), date(t.prazo), date(t.publicacao_prevista), usuarios.find(u => u.id === t.aprovador_id)?.nome, date(t.aprovado_em), t.link_publicado])].map(row => row.map(cell).join(';')).join('\r\n');
}
