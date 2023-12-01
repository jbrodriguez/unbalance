import React from 'react';

interface Props {
  steps: number;
  currentStep: number;
  labels?: { title: string; subtitle: string }[];
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

const getStatus = (currentStep: number, step: number) =>
  currentStep === step ? 'on' : 'off';

export const Stepper: React.FC<Props> = ({
  steps = 1,
  currentStep = 1,
  labels = [{ title: 'Step 1 title', subtitle: 'Step 1 subtitle' }],
}) => {
  return (
    <div>
      <ol className="items-center w-full space-y-4 sm:flex sm:space-x-8 sm:space-y-0 rtl:space-x-reverse">
        <li
          className={`flex items-center space-x-2.5 rtl:space-x-reverse ${
            styles[getStatus(currentStep, 1)].li
          }`}
        >
          <span
            className={`flex items-center justify-center w-8 h-8 border rounded-full shrink-0 ${
              styles[getStatus(currentStep, 1)].span
            }`}
          >
            1
          </span>
          <span>
            <h3 className="font-medium leading-tight">{labels[0].title}</h3>
            <p className="text-sm">{labels[0].subtitle}</p>
          </span>
        </li>
        {steps > 1 && labels.length > 1 && (
          <li
            className={`flex items-center space-x-2.5 rtl:space-x-reverse ${
              styles[getStatus(currentStep, 2)].li
            }`}
          >
            <span
              className={`flex items-center justify-center w-8 h-8 border rounded-full shrink-0 ${
                styles[getStatus(currentStep, 2)].span
              }`}
            >
              2
            </span>
            <span>
              <h3 className="font-medium leading-tight">{labels[1].title}</h3>
              <p className="text-sm">{labels[1].subtitle}</p>
            </span>
          </li>
        )}
        {steps > 2 && labels.length > 2 && (
          <li
            className={`flex items-center space-x-2.5 rtl:space-x-reverse ${
              styles[getStatus(currentStep, 3)].li
            }`}
          >
            <span
              className={`flex items-center justify-center w-8 h-8 border  rounded-full shrink-0 ${
                styles[getStatus(currentStep, 3)].span
              }`}
            >
              3
            </span>
            <span>
              <h3 className="font-medium leading-tight">{labels[2].title}</h3>
              <p className="text-sm">{labels[2].subtitle}</p>
            </span>
          </li>
        )}
      </ol>
    </div>
  );
};
