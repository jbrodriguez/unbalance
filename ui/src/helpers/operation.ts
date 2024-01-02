import { Op } from '../types';

export const operationKindToName: Record<number, string> = {
  [Op.ScatterMove]: 'Scatter / Move',
  [Op.ScatterCopy]: 'Scatter / Copy',
  [Op.ScatterValidate]: 'Scatter / Validate',
  [Op.GatherMove]: 'Gather / Move',
};
