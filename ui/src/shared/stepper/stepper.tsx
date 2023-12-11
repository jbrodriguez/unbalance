import React from 'react';

import { NavLink } from 'react-router-dom';

interface Props {
  steps: number;
  currentStep: number;
  config?: { navTo: string; title: string; subtitle: string }[];
}

const styles = {
  on: {
    li: 'text-blue-600 dark:text-blue-500',
    span: 'border-blue-600 dark:border-blue-500',
  },
  off: {
    li: 'text-gray-500 dark:text-gray-400',
    span: 'border-gray-500 dark:border-gray-400',
  },
};

const getStyles = (isActive: boolean) => (isActive ? 'on' : 'off');

export const Stepper: React.FC<Props> = ({
  steps = 1,
  config = [
    { navTo: 'select', title: 'Step 1 title', subtitle: 'Step 1 subtitle' },
  ],
}) => {
  return (
    <div>
      <ol className="items-center w-full space-y-4 sm:flex sm:space-x-8 sm:space-y-0 rtl:space-x-reverse">
        <NavLink to={config[0].navTo}>
          {({ isActive }) => (
            <li
              className={`flex items-center space-x-2.5 rtl:space-x-reverse ${
                styles[getStyles(isActive)].li
              }`}
            >
              <span
                className={`flex items-center justify-center w-8 h-8 border rounded-full shrink-0 ${
                  styles[getStyles(isActive)].span
                }`}
              >
                1
              </span>
              <span>
                <h3 className="font-medium leading-tight">{config[0].title}</h3>
                <p className="text-sm">{config[0].subtitle}</p>
              </span>
            </li>
          )}
        </NavLink>

        {steps > 1 && config.length > 1 && (
          <NavLink to={config[1].navTo}>
            {({ isActive }) => (
              <li
                className={`flex items-center space-x-2.5 rtl:space-x-reverse ${
                  styles[getStyles(isActive)].li
                }`}
              >
                <span
                  className={`flex items-center justify-center w-8 h-8 border rounded-full shrink-0 ${
                    styles[getStyles(isActive)].span
                  }`}
                >
                  2
                </span>
                <span>
                  <h3 className="font-medium leading-tight">
                    {config[1].title}
                  </h3>
                  <p className="text-sm">{config[1].subtitle}</p>
                </span>
              </li>
            )}
          </NavLink>
        )}

        {steps > 2 && config.length > 2 && (
          <NavLink to={config[2].navTo}>
            {({ isActive }) => (
              <li
                className={`flex items-center space-x-2.5 rtl:space-x-reverse ${
                  styles[getStyles(isActive)].li
                }`}
              >
                <span
                  className={`flex items-center justify-center w-8 h-8 border  rounded-full shrink-0 ${
                    styles[getStyles(isActive)].span
                  }`}
                >
                  3
                </span>
                <span>
                  <h3 className="font-medium leading-tight">
                    {config[2].title}
                  </h3>
                  <p className="text-sm">{config[2].subtitle}</p>
                </span>
              </li>
            )}
          </NavLink>
        )}
      </ol>
    </div>
  );
};
