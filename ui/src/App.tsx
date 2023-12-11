import React from 'react';

import { Outlet, useNavigate } from 'react-router-dom';

import { Header } from './shared/header/header';
import { Footer } from './shared/footer/footer';
import { useConfigActions } from './state/config';
import { useUnraidStore } from './state/unraid';

export function App() {
  const { getConfig } = useConfigActions();
  const { array } = useUnraidStore();
  const navigate = useNavigate();

  React.useEffect(() => {
    getConfig();
  }, [getConfig]);

  React.useEffect(() => {
    if (array === 'stopped') {
      navigate('/scatter');
    }
  }, [navigate, array]);

  if (array === 'stopped') {
    return null;
  }

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
