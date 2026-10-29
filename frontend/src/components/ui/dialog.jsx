import * as DialogPrimitive from '@radix-ui/react-dialog';
import { X } from 'lucide-react';
export function Dialog({ open, onOpenChange, title, children, wide = false }) {
  return <DialogPrimitive.Root open={open} onOpenChange={onOpenChange}><DialogPrimitive.Portal><DialogPrimitive.Overlay className="modal-overlay"/><DialogPrimitive.Content className={'modal ' + (wide ? 'wide-modal' : '')}><div className="modal-heading"><DialogPrimitive.Title>{title}</DialogPrimitive.Title><DialogPrimitive.Close className="icon-button" aria-label="Fechar"><X size={20}/></DialogPrimitive.Close></div><DialogPrimitive.Description className="sr-only">Gerencie os dados e acompanhe esta tarefa.</DialogPrimitive.Description>{children}</DialogPrimitive.Content></DialogPrimitive.Portal></DialogPrimitive.Root>;
}
