export const tipoLabels = { video: 'Vídeo', materia: 'Matéria', arte: 'Arte / peça gráfica', outro: 'Outra demanda' };
export const etapaLabels = { pauta: 'Pauta', producao: 'Produção', edicao: 'Edição', revisao: 'Revisão', aprovacao: 'Aprovação', pronto_publicacao: 'Pronto para publicar', publicado: 'Publicado', encerrada: 'Encerrada (legado)' };
export const editorialStages = Object.keys(etapaLabels).filter(s => s !== 'encerrada');
export const canais = ['Portal institucional', 'Instagram', 'YouTube', 'Facebook', 'Imprensa', 'Comunicação interna'];
const transitions = { pauta: ['producao'], producao: ['pauta', 'edicao', 'revisao'], edicao: ['producao', 'revisao'], revisao: ['producao', 'edicao', 'aprovacao'], aprovacao: ['revisao', 'edicao', 'pronto_publicacao'], pronto_publicacao: ['revisao', 'publicado'], publicado: ['revisao'], encerrada: ['pauta'] };
const contentFields = ['titulo', 'descricao', 'tipo_conteudo', 'canal', 'link_material', 'setor_id'];
export const statusForStage = etapa => etapa === 'pauta' ? 'pendente' : ['publicado', 'encerrada'].includes(etapa) ? 'concluida' : 'em_andamento';
export const deadline = task => task.prazo_etapa || task.prazo;
export const publicationLate = (task, now = new Date()) => !['publicado', 'encerrada'].includes(task.etapa) && task.publicacao_prevista && new Date(task.publicacao_prevista) < now;
export function validLink(link) { if (!link) return true; if (link.length > 2048) return false; try { const url = new URL(link); return ['https:', 'http:'].includes(url.protocol) && !url.username && !url.password; } catch { return false; } }
export function availableStages(task, user) {
  return [task.etapa, ...(transitions[task.etapa] || []).filter(s => !(task.tipo_conteudo === 'video' && task.etapa === 'producao' && s === 'revisao') && !(s === 'pronto_publicacao' && (user.perfil !== 'gestor' || user.setor_id !== task.setor_id)) && !(task.etapa === 'revisao' && s === 'aprovacao' && task.revisor_id && user.id !== task.revisor_id && user.perfil !== 'gestor') && !(['publicado', 'encerrada'].includes(task.etapa) && user.perfil !== 'gestor'))];
}
export function applyEditorial(previous, input, user, now = new Date()) {
  const t = { ...input };
  t.tipo_conteudo ||= 'outro'; t.etapa ||= previous?.etapa || 'pauta'; t.canal ||= 'Comunicação interna';
  if (!tipoLabels[t.tipo_conteudo] || !etapaLabels[t.etapa] || !canais.includes(t.canal)) throw new Error('Tipo, etapa ou canal editorial inválido.');
  if (!validLink(t.link_material) || !validLink(t.link_publicado)) throw new Error('Use links HTTP ou HTTPS válidos, sem credenciais.');
  if (!t.prazo) throw new Error('Informe o prazo de entrega.');
  for (const value of [t.prazo, t.prazo_etapa, t.publicacao_prevista].filter(Boolean)) if (!Number.isFinite(+new Date(value))) throw new Error('Prazo inválido.');
  if (t.prazo_etapa && new Date(t.prazo_etapa) > new Date(t.prazo)) throw new Error('O prazo da etapa deve ser anterior ou igual ao prazo de entrega.');
  if (t.publicacao_prevista && new Date(t.prazo) > new Date(t.publicacao_prevista)) throw new Error('A entrega deve acontecer até a publicação prevista.');
  t.aprovado_em = previous?.aprovado_em || null; t.aprovador_id = previous?.aprovador_id || null; t.publicado_em = previous?.publicado_em || null;
  if (!previous) {
    if (t.etapa !== 'pauta') throw new Error('Novas demandas começam na pauta.');
    t.versao = 1; t.link_publicado = '';
  } else {
    if (input.versao !== previous.versao) throw new Error('Esta demanda mudou. Reabra os detalhes antes de salvar.');
    const changed = contentFields.some(f => (t[f] || '') !== (previous[f] || ''));
    if (t.etapa !== previous.etapa && !(transitions[previous.etapa] || []).includes(t.etapa)) throw new Error('Siga as etapas de produção, revisão e aprovação antes de publicar.');
    if (previous.tipo_conteudo === 'video' && previous.etapa === 'producao' && t.etapa === 'revisao') throw new Error('O vídeo precisa passar pela edição antes da revisão.');
    if (previous.etapa === 'revisao' && t.etapa === 'aprovacao' && previous.revisor_id && user.id !== previous.revisor_id && user.perfil !== 'gestor') throw new Error('A revisão deve ser concluída pelo revisor responsável ou gestor.');
    if (['publicado', 'encerrada'].includes(previous.etapa) && (t.etapa !== previous.etapa || changed) && user.perfil !== 'gestor') throw new Error('Somente o gestor pode reabrir uma entrega finalizada.');
    if ((previous.etapa === 'aprovacao' && ['revisao', 'edicao'].includes(t.etapa)) || (['publicado', 'encerrada'].includes(previous.etapa) && (t.etapa !== previous.etapa || changed))) if (!t.anotacao?.trim()) throw new Error('Registre o motivo dos ajustes ou da reabertura no histórico.');
    if (previous.etapa === 'aprovacao' && t.etapa === 'pronto_publicacao') {
      if (user.perfil !== 'gestor' || user.setor_id !== previous.setor_id) throw new Error('A aprovação final é feita por um gestor do setor responsável.');
      if (changed) throw new Error('O conteúdo mudou durante a aprovação. Envie a nova versão para revisão.');
      t.aprovado_em = now.toISOString(); t.aprovador_id = user.id;
    }
    if (changed && (previous.aprovado_em || previous.etapa === 'aprovacao')) {
      if (t.etapa === 'publicado' && previous.etapa !== 'publicado') throw new Error('O conteúdo alterado precisa de nova revisão e aprovação.');
      t.etapa = 'revisao';
    }
    if (!['pronto_publicacao', 'publicado'].includes(t.etapa)) { t.aprovado_em = null; t.aprovador_id = null; t.publicado_em = null; t.link_publicado = ''; }
    if (t.etapa === 'publicado' && previous.etapa !== 'publicado') {
      if (!previous.aprovado_em || !previous.aprovador_id || changed) throw new Error('A publicação exige uma aprovação válida da versão atual.');
      if (!t.link_publicado?.trim()) throw new Error('Informe o link da publicação realizada.');
      t.publicado_em = now.toISOString();
    }
    if (t.etapa === 'publicado' && !t.link_publicado?.trim()) throw new Error('Mantenha o link da publicação realizada.');
    t.versao = previous.versao + 1;
  }
  t.status = statusForStage(t.etapa);
  t.concluido_em = t.status === 'concluida' ? previous?.concluido_em || now.toISOString() : null;
  return t;
}
