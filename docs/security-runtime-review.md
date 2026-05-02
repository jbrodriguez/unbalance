# Security runtime review: root process on Unraid

Branch: `security/root-runtime-review`

## Context

`unbalanced` currently performs filesystem scans and transfers across `/mnt/*` disks/shares. On Unraid this often means running with broad privileges, sometimes as `root`, because ownership and permissions across user shares/disks can block a plain `nobody` process.

The security goal should not be "never root" at all costs. The practical goal is: **keep the web/API/control plane as low-privilege as possible, and make every privileged filesystem action narrow, validated, auditable, and intentional.**

## Current high-risk architecture points

### 1. Web server and transfer executor live in the same privileged process

The Echo HTTP server, websocket command handler, planner, rsync launcher, source deletion, history, auth/session handling, and config writes all run inside the same binary/process.

If that process is running as `root`, then a bug in any HTTP/auth/websocket/config path becomes a root-level file operation risk.

Relevant code:

- `daemon/services/server/server.go` starts the HTTP server and registers the API/websocket endpoints.
- `daemon/services/core/core.go` consumes websocket packets and starts scatter/gather/move/replay/remove operations.
- `daemon/services/core/operation.go` executes rsync and source deletion.

### 2. Websocket command channel is powerful

The websocket accepts command packets and publishes them directly into the core mailbox. Auth/CSRF/origin checks exist when auth is enabled, which is good, but the command payload itself is trusted after decoding.

Risk: if an attacker gets a session, if auth is disabled, or if a future endpoint bypasses checks, the websocket can trigger transfer/delete operations.

Concrete sensitive commands:

- `scatter:move`
- `gather:move`
- `remove:source`
- `replay`

### 3. Client-supplied operation/command replay can cross trust boundaries

Some command paths accept a full `domain.Operation` or `domain.Command` from the client/websocket and use its fields later for rsync/deletion.

Relevant code:

- `CommandScatterValidate` accepts `domain.Operation` and reuses `original.Commands`.
- `CommandRemoveSource` accepts `Operation` + `Command` from the client.
- `CommandReplay` accepts `domain.Operation` from the client.

That means fields like `Src`, `Dst`, `Entry`, `RsyncArgs`, and `OpKind` should be treated as untrusted. Today they are not strongly revalidated against a server-created plan/history entry before execution.

### 4. Source deletion uses shell commands

`handleItemDeletion` builds shell commands:

- `rm -rf "<source/entry>"`
- `find "<parent>" -type d -empty -prune -exec rm -rf {} \;`

Quoting helps but does not fully remove shell risk. More importantly, `rm -rf` is high-impact when running as root. This should move to Go filesystem calls with strict path validation and symlink/path-boundary checks.

### 5. Planner paths are only partially normalized

API tree/locate paths are normalized under `/mnt`, which is good, but operation paths later come from plan/operation structures and should be rechecked at execution time.

A safe execution layer should require:

- `Src` resolves under `/mnt/disk*` or another explicitly allowed Unraid disk root.
- `Dst` resolves under `/mnt/disk*`.
- `Entry` is relative, clean, non-empty, not `.` or `..`, and does not escape when joined.
- resolved source/destination are not symlinks escaping the allowed roots.
- source deletion only touches the exact source path that was transferred.

### 6. `rsyncArgs` is user-configurable and flows into execution

`RsyncArgs` is set from config and appended into rsync execution. It is passed as argv, not shell-concatenated, which is good. Still, because rsync supports dangerous options (`--delete`, `--remove-source-files`, `--rsync-path`, remote syntax, etc.), the app should validate or constrain custom args.

This is especially important if the process runs as root.

### 7. Global CORS is open by default

`middleware.CORS()` is enabled with defaults. CSRF and SameSite help for authenticated mutations, and websocket origin is checked when auth is enabled. Still, the app should prefer a same-origin CORS policy by default, especially for a local root-capable service.

### 8. Default auth can be disabled

Auth support has improved, but `AUTH_ENABLED=false` remains the default. In a trusted LAN appliance this may be acceptable historically, but for a root-capable file mover, the safer posture is:

- make auth strongly recommended in UI/docs;
- warn loudly when disabled;
- require explicit opt-out for dangerous operations when unauthenticated; or
- consider enabling auth by default for new installs while preserving upgrades.

## Best improvement opportunities

### P1: Introduce a server-side operation capability model

Do not let client-submitted `Operation` / `Command` objects directly drive execution.

Recommended shape:

1. Planning creates a server-side `planID` and stores the immutable normalized command list in memory/history.
2. The UI can request `execute(planID)` / `validate(planID)` / `removeSource(commandID)`.
3. The server looks up the original plan/command and executes only that stored version.
4. Any replay from history should rehydrate a stored server history item, not trust client-sent paths/args.

This is the cleanest architectural security win.

### P1: Add a path safety package and enforce it at execution time

Create a small package/function set responsible for validating execution paths before every rsync/delete:

- allowed roots from current Unraid disk discovery, not hardcoded strings only;
- clean relative entries;
- no absolute `Entry`;
- no `..` escape;
- symlink escape detection where possible;
- source/destination must be different roots;
- deletion must validate again immediately before removal.

This should be applied in `runCommand`, `handleItemDeletion`, `removeSource`, `replay`, and validate paths.

### P1: Replace shell-based deletion with safe Go deletion primitives

Replace `rm -rf` shell execution with guarded Go code:

- `os.Remove` for files;
- `os.RemoveAll` only after safe path validation and maybe refusing root-ish paths;
- parent pruning via walking upward and `os.Remove` only empty directories, stopping before share/disk boundary.

This removes shell injection risk and gives better guardrails/logging.

### P2: Split low-privilege UI/API from privileged worker

Longer-term architecture:

- web/API process runs as `nobody`;
- a small local privileged worker performs only validated filesystem actions;
- communication over a Unix socket with file permissions, or a very small localhost-only RPC;
- commands are structured, not shell strings;
- worker owns path validation and audit logging.

This is more work, but it materially reduces blast radius. Even if the worker remains root, the attack surface that is root-capable becomes much smaller.

### P2: Constrain `rsyncArgs`

Adopt an allowlist/denylist policy for custom rsync args.

At minimum, block or warn on:

- remote execution/transport options;
- `--delete*` unless explicitly supported;
- `--remove-source-files`;
- options that invoke external programs;
- options that mutate ownership/perms beyond expected behavior.

Also store parsed args as structured tokens; avoid round-tripping through `RsyncStrArgs` except for display.

### P2: Harden HTTP defaults

- Replace default open CORS with same-origin/no CORS unless dev mode is enabled.
- Add security headers (`X-Content-Type-Options`, `Referrer-Policy`, conservative CSP where compatible).
- Warn in UI when auth is disabled.
- Consider shorter sessions than 180 days or configurable session duration.
- Ensure session files and env files are written with restrictive permissions where Unraid permits it.

### P3: Privilege drop / staged privilege approach

If Unraid allows it, start as root only to bind/read what is needed, then drop privileges for the web process. For transfers that need elevated permissions, use one of:

- a root worker;
- controlled `sudo` rules for rsync/delete only;
- Linux capabilities where applicable, although file ownership/permission bypass usually still needs broad privileges.

Given Unraid constraints, the worker split is probably more realistic than trying to perfect `nobody` permissions globally.

## Suggested first implementation path

1. Add `daemon/services/core/safepath.go` with validation helpers and tests.
2. Enforce safe paths in `createScatterOperation`, `createGatherOperation`, `runCommand`, and `handleItemDeletion`.
3. Replace shell deletion with validated Go deletion/pruning.
4. Change websocket command payloads to use IDs rather than full executable `Operation` objects for execute/replay/remove-source.
5. Then evaluate a privileged-worker split once the trust boundary is explicit.
