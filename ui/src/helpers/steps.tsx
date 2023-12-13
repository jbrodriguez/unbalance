import { Op, Step } from '~/types';

export function convertStatusToStep(status: Op): Step {
  switch (status) {
    case Op.Neutral:
      return 'select';
    case Op.ScatterPlan:
    case Op.GatherPlan:
      return 'plan';
    case Op.ScatterMove:
    case Op.GatherMove:
    case Op.ScatterCopy:
    case Op.ScatterValidate:
      return 'transfer';
    default:
      return 'select';
  }
}

export function getRouteFromOp(status: Op): string {
  switch (status) {
    case Op.Neutral:
      return '/scatter/select';
    case Op.ScatterPlan:
      return '/scatter/plan';
    case Op.GatherPlan:
      return '/gather/plan';
    case Op.ScatterMove:
    case Op.ScatterCopy:
    case Op.ScatterValidate:
      return '/scatter/transfer';
    case Op.GatherMove:
      return '/gather/transfer';
    default:
      return '/scatter/select';
  }
}

// export function getRoute(step: Step): string {
//   switch (step) {
//     case 'idle':
//       return '/scatter';
//     case 'scatter.select':
//       return '/scatter/select';
//     case 'scatter.plan':
//       return '/scatter/plan';
//     case 'scatter.transfer':
//       return '/scatter/transfer';
//     case 'gather.select':
//       return '/gather/select';
//     case 'gather.plan':
//       return '/gather/plan';
//     case 'gather.transfer':
//       return '/gather/transfer';
//     case 'history':
//       return '/history';
//     case 'settings.notications':
//       return '/settings/notifications';
//     case 'settings.reserved':
//       return '/settings/reserved';
//     case 'settings.rsync':
//       return '/settings/rsync';
//     case 'settings.verbosity':
//       return '/settings/verbosity';
//     case 'settings.update':
//       return '/settings/update';
//     case 'settings.refresh':
//       return '/settings/refresh';
//     case 'logs':
//       return '/logs';
//     default:
//       return '/scatter';
//   }
// }
