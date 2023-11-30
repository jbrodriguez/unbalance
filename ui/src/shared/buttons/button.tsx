import React from 'react';

interface Props {
  label: string;
  variant?: 'primary' | 'secondary' | 'accent';
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  onClick?: () => void;
  disabled?: boolean;
}

const variants = {
  primary:
    'text-white bg-blue-700 hover:bg-blue-800 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700  dark:focus:ring-blue-800',
  secondary:
    'text-gray-900 bg-white border border-gray-300 hover:bg-gray-100 focus:ring-gray-200 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700',
  accent:
    'text-white bg-green-700 hover:bg-green-800 focus:ring-green-300 dark:bg-green-600 dark:hover:bg-green-700 dark:focus:ring-green-800',
};

export const Button: React.FC<Props> = ({
  label = 'Ok',
  variant = 'primary',
  leftIcon = null,
  rightIcon = null,
  onClick = () => {},
  disabled = false,
}) => {
  const disabledStyle = disabled ? 'opacity-50 cursor-not-allowed' : '';
  const variantStyle = variants[variant];

  return (
    <button
      type="button"
      className={`flex flex-row items-center font-medium rounded-lg text-sm px-5 py-2.5 focus:ring-4 focus:outline-none ${variantStyle} ${disabledStyle}`}
      onClick={onClick}
      disabled={disabled}
    >
      {leftIcon}
      {label}
      {rightIcon}
    </button>
  );
};
