import React from 'react';

import { Input } from '@/components/ui/input';

import { Panel } from '~/shared/panel/panel';
import { useUnraidActions, useUnraidLogs } from '~/state/unraid';
import { parseLogLine, LogLevel } from '~/helpers/logs';

type Filter = 'all' | LogLevel;

const levelStyles: Record<LogLevel, string> = {
  info: 'text-gray-700 dark:text-gray-500',
  warn: 'text-yellow-600 dark:text-yellow-500',
  error: 'text-red-600 dark:text-red-500',
};

const filters: Array<Filter> = ['all', 'info', 'warn', 'error'];

export const Logs: React.FunctionComponent = () => {
  const { getLog } = useUnraidActions();
  const logs = useUnraidLogs();
  const [filter, setFilter] = React.useState<Filter>('all');
  const [search, setSearch] = React.useState('');

  React.useEffect(() => {
    getLog();
  }, [getLog]);

  const parsed = logs.map(parseLogLine);
  const counts = parsed.reduce(
    (acc, line) => {
      acc[line.level] += 1;
      return acc;
    },
    { info: 0, warn: 0, error: 0 },
  );

  const term = search.toLowerCase();
  const visible = parsed.filter(
    (line) =>
      (filter === 'all' || line.level === filter) &&
      (term === '' || line.message.toLowerCase().includes(term)),
  );

  const count = (value: Filter) =>
    value === 'all' ? parsed.length : counts[value];

  const subtitle = (
    <div className="flex flex-1 flex-row items-center gap-2">
      {filters.map((value) => (
        <button
          key={value}
          onClick={() => setFilter(value)}
          className={
            filter === value
              ? 'px-2 py-0.5 text-xs rounded bg-blue-700 text-white'
              : 'px-2 py-0.5 text-xs rounded bg-gray-200 text-gray-700 dark:bg-gray-800 dark:text-gray-400 hover:bg-gray-300 dark:hover:bg-gray-700'
          }
        >
          {value} ({count(value)})
        </button>
      ))}
      <div className="flex-1" />
      <Input
        className="h-7 w-64"
        placeholder="filter messages ..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />
    </div>
  );

  return (
    <Panel title="Logs" subtitle={subtitle}>
      <div className="font-mono text-sm">
        {visible.map((line, index) => (
          <p key={index} className={levelStyles[line.level]}>
            {line.timestamp && (
              <span className="text-gray-400 dark:text-gray-600 pr-2">
                {line.timestamp}
              </span>
            )}
            {line.message}
          </p>
        ))}
        {visible.length === 0 && (
          <p className="text-gray-500">no log lines match the current filter</p>
        )}
      </div>
    </Panel>
  );
};
