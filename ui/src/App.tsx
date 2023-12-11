import React from 'react';

import { Outlet } from 'react-router-dom';

import { Header } from './shared/header/header';
import { Footer } from './shared/footer/footer';
import { useConfigActions } from './state/config';

export function App() {
  const { getConfig } = useConfigActions();

  React.useEffect(() => {
    getConfig();
  }, [getConfig]);

  return (
    <div className="container mx-auto h-screen flex flex-col">
      <header>
        <Header />
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
      <footer>
        <Footer />
      </footer>
    </div>
  );
}
