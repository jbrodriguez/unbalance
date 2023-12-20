import { Op } from '~/types';

// export function convertStatusToStep(status: Op): Step {
//   switch (status) {
//     case Op.Neutral:
//       return 'select';
//     case Op.ScatterPlan:
//     case Op.GatherPlan:
//       return 'plan';
//     case Op.ScatterMove:
//     case Op.GatherMove:
//     case Op.ScatterCopy:
//     case Op.ScatterValidate:
//       return 'transfer';
//     default:
//       return 'select';
//   }
// }

export function getRouteFromStatus(status: Op): string {
  switch (status) {
    case Op.Neutral:
      return '/scatter/select';
    case Op.ScatterPlanning:
      return '/scatter/plan/log';
    case Op.ScatterPlan:
      return '/scatter/plan/validation';
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

// export const stepToIndex = {
//   idle: 1,
//   select: 1,
//   plan: 2,
//   transfer: 3,
// };

export const routeToStep = (route: string): number => {
  switch (route) {
    case '/':
    case '/scatter':
    case '/scatter/select':
    case '/gather':
    case '/gather/select':
      return 1;
    case '/scatter/plan':
    case '/scatter/plan/log':
    case '/scatter/plan/validation':
    case '/gather/plan':
      return 2;
    case '/scatter/transfer':
    case '/gather/transfer':
      return 3;
    default:
      return 1;
  }
};

export const getNextRoute = (route: string) => {
  switch (route) {
    case '/':
      return '/scatter/select';
    case '/scatter':
      return '/scatter/select';
    case '/scatter/select':
      return '/scatter/plan/log';
    case '/scatter/plan':
    case '/scatter/plan/log':
    case '/scatter/plan/validation':
      return '/scatter/transfer';
    case '/gather':
      return '/gather/select';
    case '/gather/select':
      return '/gather/plan';
    case '/gather/plan':
      return '/gather/transfer';
    case '/history':
      return '/history';
    case '/settings':
      return '/settings';
    case '/logs':
      return '/logs';
    default:
      return '/scatter/select';
  }
};

export const getBaseRoute = (path: string) => {
  // Split the path using '/' as the delimiter
  const parts = path.split('/');

  // Extract the first non-empty part
  const firstLevel = parts.find((part) => part !== '');

  // Return the first level with the leading '/'
  return '/' + firstLevel;
};
