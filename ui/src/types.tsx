export enum Op {
  Neutral = 0,
  ScatterPlanning = 1,
  ScatterPlan = 2,
  ScatterMove = 3,
  ScatterCopy = 4,
  ScatterValidate = 5,
  GatherPlanning = 6,
  GatherPlan = 7,
  GatherMove = 8,
}

export type Step = 'idle' | 'select' | 'plan' | 'transfer';

export type Variant = 'primary' | 'secondary' | 'accent';

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
  currentFree: number;
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
  vdisks: { [key: string]: VDisk };
  bytesToTransfer: number;
}

export interface State {
  status: number;
  unraid: Unraid | null;
  operation: Operation | null;
  history: History | null;
  plan: Plan | null;
}

export interface Node {
  id: string;
  label: string;
  leaf: boolean;
  parent: string;
  checked?: boolean;
  expanded?: boolean;
  loading?: boolean;
  children: string[];
}

export type Nodes = Record<string, Node>;

export interface Icons {
  collapseIcon: React.ReactElement;
  expandIcon: React.ReactElement;
  checkedIcon: React.ReactElement;
  uncheckedIcon: React.ReactElement;
  parentIcon: React.ReactElement;
  leafIcon: React.ReactElement;
  hiddenIcon: React.ReactElement;
  loadingIcon: React.ReactElement;
}

export interface Branch {
  nodes: Nodes;
  order: string[];
}

export type Chosen = Record<string, boolean>;
export type Targets = Record<string, boolean>;

export enum Topic {
  CommandScatterPlanStart = 'scatter:plan:start',
  EventScatterPlanStarted = 'scatter:plan:started',
  EventScatterPlanProgress = 'scatter:plan:progress',
  EventScatterPlanEnded = 'scatter:plan:ended',
}

export interface Packet {
  topic: Topic;
  payload: string;
}
