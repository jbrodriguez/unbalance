# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Go Backend
- `make release` - Build production binary for Linux (cross-compiled from macOS)
- `make local` - Build development binary for current platform
- `make test` - Run Go unit tests
- `make clean` - Clean Go build artifacts
- `make kill` - Stop running unbalanced daemon

### React Frontend (ui/ directory)
- `cd ui && npm run dev` - Start development server with hot reload
- `cd ui && npm run build` - Build production frontend assets
- `cd ui && npm run lint` - Run ESLint on TypeScript/React code
- `cd ui && npm run preview` - Preview production build

### Full Development Workflow
For frontend development with live backend proxy:
1. Modify `ui/vite.config.ts` to add proxy configuration pointing to your Unraid server
2. Run `npm run dev` in ui/ directory
3. Access localhost:5173 for development

## Architecture Overview

**unbalanced** is an Unraid plugin for transferring files between disks in arrays. It's a hybrid Go/React application with embedded frontend.

### Backend Structure (Go)
- **Entry Point**: `unbalance.go` - CLI parsing and application bootstrap
- **Core Services**: `daemon/services/core/` - Business logic for file operations
- **Domain Models**: `daemon/domain/` - Data structures and types
- **REST API**: `daemon/services/server/` - Echo web server with WebSocket support
- **Algorithms**: `daemon/algorithm/` - Greedy and knapsack algorithms for space optimization

### Frontend Structure (React/TypeScript)
- **Main App**: `ui/src/App.tsx` - Root component with routing
- **Flows**: `ui/src/flows/` - Main application workflows (gather, scatter, history, settings, logs)
- **Shared Components**: `ui/src/shared/` - Reusable UI components
- **State Management**: `ui/src/state/` - Zustand stores for application state
- **API Layer**: `ui/src/api/` - HTTP client and WebSocket communication

### Key Architectural Patterns
- **Pub/Sub Communication**: Uses `github.com/cskr/pubsub` for internal event messaging
- **WebSocket Real-time Updates**: Live progress updates during file operations
- **Embedded Assets**: Frontend built into Go binary using embed directive
- **Two Operation Modes**:
  - **Scatter**: Transfer files from one disk to multiple target disks
  - **Gather**: Consolidate files from multiple locations into single disk

### File Operations
- Uses rsync for actual file transfers with customizable flags
- Default rsync args: `-X` (preserve extended attributes)
- Operations support both MOVE and COPY modes
- Built-in validation and dry-run capabilities

### Configuration
- Environment variables and CLI flags for runtime configuration
- Settings persisted and configurable via web UI
- Supports custom rsync arguments for power users

### Plugin Integration
Designed specifically for Unraid systems:
- Reads from `/mnt/user/` for user shares
- Operates on `/mnt/disk*/` paths for individual disks
- Integrates with Unraid notification system