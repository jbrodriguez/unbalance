// import React from 'react';

import AutoSizer from 'react-virtualized-auto-sizer';

import { CheckboxTree } from '~/shared/tree/checkbox-tree';
import { Node } from '~/types';
import { Icon } from '~/shared/icons/icon';
import { useScatterTree, useScatterActions } from '~/state/scatter';

interface Props {
  height?: number;
  width?: number;
}

export const FileSystem: React.FunctionComponent<Props> = ({ height }) => {
  const tree = useScatterTree();
  const { loadBranch, toggleSelected } = useScatterActions();

  const onLoad = async (node: Node) => await loadBranch(node);
  const onCheck = (node: Node) => toggleSelected(node);

  return (
    <div className="flex flex-1 bg-neutral-200 dark:bg-gray-950">
      <AutoSizer disableHeight>
        {({ width }) => (
          <div
            className="overflow-y-auto overflow-x-auto p-2"
            style={{ height: `${height}px`, width: `${width}px` }}
          >
            <CheckboxTree
              nodes={tree}
              onLoad={onLoad}
              onCheck={onCheck}
              icons={{
                collapseIcon: (
                  <Icon
                    name="minus"
                    size={20}
                    fill="fill-slate-500 dark:fill-gray-700"
                  />
                ),
                expandIcon: (
                  <Icon
                    name="plus"
                    size={20}
                    fill="fill-slate-500 dark:fill-gray-700"
                  />
                ),
                checkedIcon: (
                  <Icon
                    name="checked"
                    size={20}
                    fill="fill-slate-700 dark:fill-lime-600"
                  />
                ),
                uncheckedIcon: (
                  <Icon
                    name="unchecked"
                    size={20}
                    fill="fill-slate-700 dark:fill-slate-200"
                  />
                ),
                leafIcon: (
                  <Icon
                    name="file"
                    size={20}
                    fill="fill-blue-400 dark:fill-gray-700"
                  />
                ),
                parentIcon: (
                  <Icon
                    name="folder"
                    size={20}
                    fill="fill-orange-400 dark:fill-gray-700"
                  />
                ),
                hiddenIcon: (
                  <Icon
                    name="square"
                    size={20}
                    fill="fill-neutral-200 dark:fill-gray-950"
                  />
                ),
                loadingIcon: (
                  <Icon
                    name="loading"
                    size={20}
                    fill="animate-spin fill-slate-700 dark:fill-slate-700"
                  />
                ),
              }}
            />
          </div>
        )}
      </AutoSizer>
    </div>
  );
};
