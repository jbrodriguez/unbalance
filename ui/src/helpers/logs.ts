export type LogLevel = 'info' | 'warn' | 'error';

export interface LogLine {
  timestamp: string;
  level: LogLevel;
  message: string;
}

// matches lines written by the daemon logger, e.g.
// 2026/07/16 20:15:04 [error] unable to remove source folder
const reLogLine =
  /^(\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2})(?: \[(info|warn|error)\])? (.*)$/;

// lines written before level tags existed (or by the standard logger)
// carry no tag and are treated as info
export const parseLogLine = (line: string): LogLine => {
  const match = reLogLine.exec(line);
  if (!match) {
    return { timestamp: '', level: 'info', message: line };
  }

  return {
    timestamp: match[1],
    level: (match[2] as LogLevel) ?? 'info',
    message: match[3],
  };
};
