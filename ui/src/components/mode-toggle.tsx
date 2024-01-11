// import { Moon, Sun } from 'lucide-react';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useTheme } from '@/components/use-theme';

import { Icon } from '~/shared/icons/icon';

export function ModeToggle() {
  const { setTheme } = useTheme();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="icon"
          className="dark:bg-transparent bg-transparent dark:hover:bg-gray-700 hover:bg-slate-100  "
        >
          <Icon
            name="white-balance-sunny"
            style="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all fill-amber-500 dark:-rotate-90 dark:scale-0"
          />
          <Icon
            name="weather-night"
            style="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all fill-blue-500 dark:rotate-0 dark:scale-100"
          />
          <span className="sr-only">Toggle theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setTheme('light')}>
          Light
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme('dark')}>
          Dark
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme('system')}>
          System
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
