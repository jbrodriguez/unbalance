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
    case Op.ScatterPlanning:
    case Op.ScatterPlan:
      return '/scatter/plan';
    case Op.GatherPlanning:
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

export const stepToIndex = {
  idle: 1,
  select: 1,
  plan: 2,
  transfer: 3,
};

export const routeToIndex: { [x: string]: number } = {
  '/': 1,
  '/scatter/select': 1,
  '/scatter/plan': 2,
  '/scatter/transfer': 3,
};
