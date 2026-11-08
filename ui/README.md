# unbalanced UI

React + TypeScript + Vite frontend for the unbalanced application.

## Development

```bash
npm install
npm run dev       # dev server at http://localhost:5173
npm run build     # production build
npm run lint      # lint
```

A running backend is expected at `http://localhost:7090`. The Vite dev server proxies `/api` requests there (see `vite.config.ts`).

## Responsive Design

The UI targets desktop (≥768px). Layout uses Tailwind's `md:` breakpoint as the primary switch between mobile and desktop layouts.
