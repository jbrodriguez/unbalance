import React from 'react';

import Tree, { BasicDataNode } from 'rc-tree';
import 'rc-tree/assets/index.css';

interface Props {
  height?: number;
}

export type EventDataNode<BasicDataNode> = {
  key: React.Key;
  expanded: boolean;
  selected: boolean;
  checked: boolean;
  loaded: boolean;
  loading: boolean;
  halfChecked: boolean;
  dragOver: boolean;
  dragOverGapTop: boolean;
  dragOverGapBottom: boolean;
  pos: string;
  active: boolean;
} & BasicDataNode;

const initialData = [{ title: 'mnt', key: '/mnt', isLeaf: false }];

type CheckedKeysType =
  | { checked: React.Key[]; halfChecked: React.Key[] }
  | React.Key[];

export function FileSystem({ height }: Props) {
  const [data, setData] = React.useState<BasicDataNode[]>(initialData);
  const [checked, setChecked] = React.useState<CheckedKeysType>([]);

  const onSelect = (info: unknown) => {
    console.log('selected', info);
  };

  const onCheck = (checkedKeys: CheckedKeysType) => {
    console.log(checkedKeys);
    setChecked(checkedKeys);
  };

  const getTree = async (treeNode: EventDataNode<BasicDataNode>) => {
    // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
    treeNode.children = [{ title: 'films', key: '/mnt/films', isLeaf: false }];
    return treeNode;
  };

  const onLoadData = async (treeNode: EventDataNode<BasicDataNode>) => {
    console.log('load data... ', treeNode);

    const loaded = await getTree(treeNode);
    const treeData = [...data];
    // @ts-expect-error ts-migrate(2531) FIXME: Object is possibly 'null'.
    const index = treeData.findIndex((e) => e.key === treeNode.key);
    if (index > -1) {
      treeData[index] = loaded;
    }

    setData(treeData);
  };

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <div className={`overflow-y-auto`} style={{ height: `${height}px` }}>
        <Tree
          onSelect={onSelect}
          checkable
          onCheck={onCheck}
          checkedKeys={checked}
          loadData={onLoadData}
          treeData={data}
        />
      </div>
    </div>
  );
}
