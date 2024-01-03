import { Op, CommandStatus } from '../types';

import { Icon } from '~/shared/icons/icon';

export const operationKindToName: Record<number, string> = {
  [Op.ScatterMove]: 'Scatter / Move',
  [Op.ScatterCopy]: 'Scatter / Copy',
  [Op.ScatterValidate]: 'Scatter / Validate',
  [Op.GatherMove]: 'Gather / Move',
};

export const getCommandStatus = (status: CommandStatus): React.ReactNode => {
  switch (status) {
    case CommandStatus.Complete:
      return (
        <Icon
          name="check-circle"
          size={14}
          style="fill-green-600 dark:fill-green-600"
        />
      );
    case CommandStatus.Pending:
      return (
        <Icon
          name="minus-circle"
          size={14}
          style="fill-blue-600 dark:fill-blue-600"
        />
      );
    case CommandStatus.Flagged:
      return (
        <Icon
          name="check-circle"
          size={14}
          style="fill-yellow-600 dark:fill-yellow-600"
        />
      );
    case CommandStatus.Stopped:
      return (
        <Icon
          name="minus-circle"
          size={14}
          style="fill-red-600 dark:fill-red-600"
        />
      );
    case CommandStatus.SourceRemoval:
      return (
        <Icon
          name="loading"
          size={14}
          style="fill-yellow-600 dark:fill-yellow-600 animate-spin"
        />
      );
    default:
      return (
        <Icon
          name="loading"
          size={14}
          style="fill-slate-600 dark:fill-slate-600 animate-spin"
        />
      );
  }
};
