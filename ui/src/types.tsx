export enum Op {
  Neutral = 0,
  ScatterPlan = 1,
  ScatterMove = 2,
  ScatterCopy = 3,
  ScatterValidate = 4,
  GatherPlan = 5,
  GatherMove = 6,
}

export type Step = 'idle' | 'select' | 'plan' | 'transfer';

// export type Step =
//   | 'idle'
//   | 'scatter.select'
//   | 'scatter.plan'
//   | 'scatter.transfer'
//   | 'gather.select'
//   | 'gather.plan'
//   | 'gather.transfer'
//   | 'history'
//   | 'settings.notications'
//   | 'settings.reserved'
//   | 'settings.rsync'
//   | 'settings.verbosity'
//   | 'settings.update'
//   | 'settings.refresh'
//   | 'logs';

export interface Config {
  version: string;
  dryRun: boolean;
  notifyPlan: number;
  notifyTransfer: number;
  reservedAmount: bigint;
  reservedUnit: string;
  rsyncArgs: string[];
  verbosity: number;
  checkForUpdate: number;
  refreshRate: number;
}

export interface Unraid {
  numDisks: number;
  numProtected: number;
  synced: Date;
  syncErrs: number;
  resync: number;
  resyncPos: number;
  state: string;
  size: number;
  free: number;
  disks: Disk[];
  blockSize: number;
}

export interface Disk {
  id: number;
  name: string;
  path: string;
  device: string;
  type: string;
  fsType: string;
  free: number;
  size: number;
  serial: string;
  status: string;
  blocksTotal: number;
  blocksFree: number;
}

export interface Command {
  id: string;
  src: string;
  dst: string;
  entry: string;
  size: number;
  transferred: number;
  status: number;
}

export interface Operation {
  id: string;
  opKind: number;
  started: Date;
  finished: Date;
  bytesToTransfer: number;
  bytesTransferred: number;
  dryRun: boolean;
  rsyncArgs: string[];
  rsyncStrArgs: string;
  commands: Command[];
  completed: number;
  speed: number;
  remaining: string;
  deltaTransfer: number;
  line: string;
}

export interface History {
  version: number;
  lastChecked: Date;
  items: { [key: string]: Operation };
  order: string[];
}

export interface Item {
  Name: string;
  Size: number;
  Path: string;
  Location: string;
  BlocksUsed: number;
}

export interface Bin {
  Size: number;
  Items: Item[];
  BlocksUsed: number;
}

export interface VDisk {
  path: string;
  plannedFree: number;
  bin: Bin;
  src: boolean;
  dst: boolean;
}

export interface Plan {
  started: Date;
  finished: Date;
  chosenFolders: string[];
  ownerIssue: number;
  groupIssue: number;
  folderIssue: number;
  fileIssue: number;
  vDisks: { [key: string]: VDisk };
  bytesToTransfer: number;
}

export interface State {
  status: number;
  unraid: Unraid | null;
  operation: Operation | null;
  history: History | null;
  plan: Plan | null;
}
