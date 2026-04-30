# Privileged worker split design

Branch: `security/root-runtime-review`

## Goal

Reduce the blast radius of running `unbalanced` on Unraid when filesystem operations require root-like privileges.

Instead of one root process serving HTTP, handling auth/session state, parsing websocket commands, planning operations, and executing file moves/deletes, split the system into:

1. **UI/API daemon** — low privilege, ideally `nobody`.
2. **Filesystem worker** — small privileged process with a narrow command surface.

The goal is not to magically avoid privilege. The goal is to make the root-capable part tiny, boring, validated, and auditable.

## High-level architecture

```text
Browser
  |
  | HTTP/WebSocket
  v
unbalanced-web  -- Unix socket / local RPC -->  unbalanced-worker
(nobody)                                      (root or elevated)
  |                                               |
  | auth, sessions, UI, planning                  | validated stat/scan/rsync/delete only
  | operation state, history                      | no HTTP, no browser input, no config UI
  v                                               v
/boot/config/plugins/unbalanced            /mnt/disk*, /mnt/user as needed
```

## Process responsibilities

### `unbalanced-web` low-privilege process

Runs as `nobody` where possible.

Owns:

- HTTP server
- embedded UI assets
- auth/session/CSRF
- websocket connection to browser
- config UI and non-secret config state
- operation state shown to the frontend
- high-level planning flow
- user confirmation flow
- history presentation
- worker client

Does **not** directly:

- shell out to `rsync`
- remove files/folders
- walk arbitrary root-owned trees unless permissions allow
- trust browser-supplied executable `Operation` objects

### `unbalanced-worker` privileged process

Runs as root, or with the minimum privilege Unraid realistically allows.

Owns only filesystem capability:

- discover disks/shares
- stat/read directory metadata needed for planning
- calculate item sizes when low-privilege process cannot
- execute `rsync` with validated args
- remove exact transferred source paths
- prune empty parent directories within validated boundaries
- stream progress/events back to web process

Does **not** expose HTTP.
Does **not** parse browser sessions/cookies.
Does **not** accept raw shell commands.
Does **not** accept arbitrary absolute paths without validation.

## Communication channel

Recommended first version: Unix domain socket.

Example path:

```text
/var/run/unbalanced/worker.sock
```

Permissions:

```text
owner: root
 group: nobody/users or dedicated unbalanced group
 mode: 0660
```

The worker should reject requests from unexpected peers if peer credential checks are available. On Linux, Go can inspect Unix socket peer credentials via `SO_PEERCRED`.

If Unraid makes socket permissions awkward, a localhost-only TCP socket with a boot-generated shared token is possible, but Unix socket is cleaner.

## RPC shape

Keep the protocol structured and boring. JSON over Unix socket is enough initially; gRPC is probably overkill.

Suggested request model:

```go
type WorkerRequest struct {
    ID      string          `json:"id"`
    Method  string          `json:"method"`
    Payload json.RawMessage `json:"payload"`
}

type WorkerResponse struct {
    ID      string          `json:"id"`
    OK      bool            `json:"ok"`
    Error   string          `json:"error,omitempty"`
    Payload json.RawMessage `json:"payload,omitempty"`
}
```

Suggested methods:

- `StorageStatus`
- `ListTree`
- `LocatePath`
- `ScanItems`
- `PlanStats`
- `ExecuteRsync`
- `ValidateRsync`
- `RemoveTransferredSource`
- `StopOperation`

For long-running work, either:

1. request starts a job and returns `jobID`, then web subscribes to progress; or
2. one streaming connection emits progress events until completion.

Simpler first version: keep one worker connection open and stream newline-delimited JSON events.

## Command model

Avoid client-submitted executable operations.

Current-ish flow:

```text
browser sends full Operation/Command-ish payload
server executes fields from payload
```

Safer flow:

```text
browser selects folders/targets
web asks worker/web planner for plan
web stores immutable server-side plan as planID
browser confirms planID
web asks worker to execute planID's stored commands
worker validates every command again before execution
```

The browser should never be the source of truth for:

- `Src`
- `Dst`
- `Entry`
- `RsyncArgs`
- deletion target
- operation kind

## Path validation boundary

The worker owns final validation because it is the privileged side.

Every executable command should pass through a single validator:

```go
type SafeCommand struct {
    SrcDiskRoot string
    DstDiskRoot string
    Entry       string
    Args        []string
}
```

Validation rules:

- source root must be one of current discovered disk roots, e.g. `/mnt/disk1`
- destination root must be one of current discovered disk roots
- source and destination roots must differ
- `Entry` must be relative, clean, non-empty, not `.`, not `..`
- joined source path must remain inside source root
- joined destination path must remain inside destination root
- no remote rsync syntax in source/destination
- rsync args must pass allow/deny policy
- deletion target must exactly match a completed transfer source
- parent pruning must stop before disk root/share boundary
- symlink handling should refuse or explicitly resolve escapes before destructive actions

## Worker install/start model on Unraid

There are two practical options.

### Option A: one binary, two modes

Build one `unbalanced` binary with subcommands/modes:

```bash
unbalanced web --port 7090 --worker-socket /var/run/unbalanced/worker.sock
unbalanced worker --socket /var/run/unbalanced/worker.sock
```

Pros:

- less packaging complexity
- shared code/types
- easiest migration path

Cons:

- need discipline to keep worker imports/surface small

### Option B: two binaries

```bash
unbalanced-web
unbalanced-worker
```

Pros:

- stronger conceptual separation
- easier to audit what the worker can do

Cons:

- more packaging/release complexity

Recommendation: **start with Option A** using internal packages that enforce separation. If it feels clean, split binaries later.

## Start script shape

Current script starts the app through `sudo -H bash -c ...`.

A split version could look like:

```bash
# start privileged worker as root
nohup "$prog" worker \
  --socket /var/run/unbalanced/worker.sock \
  >> /var/log/unbalanced-worker.log 2>&1 &

# start low-privilege web process
nohup sudo -u nobody -H "$prog" web \
  --port "$PORT" \
  --worker-socket /var/run/unbalanced/worker.sock \
  >> /var/log/unbalanced-web.log 2>&1 &
```

Unraid specifics may require using the available `sudo`/`su`/`runuser` behavior. The important piece is: only the worker stays privileged.

## Migration strategy

### Phase 1: draw the boundary without changing processes

- Create a `filesystem.Executor` interface.
- Move rsync/delete/tree/stat calls behind that interface.
- Add central path validation.
- Keep everything in one process initially.

This gives security value immediately and makes tests possible.

### Phase 2: local in-process worker implementation

- Make core call an executor that looks like a worker client.
- Use server-side plan IDs and immutable command records.
- Stop accepting executable `Operation` objects from browser payloads.

### Phase 3: Unix socket worker

- Add `unbalanced worker` mode.
- Add `unbalanced web` mode.
- Replace in-process executor with Unix socket RPC client.
- Keep a compatibility mode for old single-process startup if needed.

### Phase 4: harden worker

- Peer credential check.
- Worker method allowlist.
- Strict rsync arg policy.
- Audit log of every privileged action.
- Optional operation confirmation token: worker executes only commands minted by web-side planner and signed/HMACed with a boot secret.

## Main tradeoffs

### Benefits

- HTTP/auth/UI bugs no longer directly equal root filesystem control.
- Privileged code surface becomes much smaller.
- Easier to audit destructive operations.
- Creates a natural place for path and rsync policy enforcement.
- Makes future sandboxing/capability work easier.

### Costs

- More process lifecycle complexity.
- More error states: worker missing, socket permission issue, worker restart during operation.
- Need progress streaming across process boundary.
- Some planning scans may still require worker privileges, so the boundary must be designed carefully.

## Open questions

- Can Unraid reliably run the web process as `nobody` from plugin scripts?
- Which user/group should own `/var/run/unbalanced/worker.sock`?
- Should planning scans run in the worker, or only execution/deletion?
- Should auth become required when privileged worker mode is enabled?
- How much backward compatibility is needed for single-process mode?

## Recommendation

Do this incrementally.

First implement the **executor boundary + server-side plan IDs + path validation** while still running as one process. That reduces risk immediately and avoids a large risky rewrite.

Then introduce the worker mode once the internal boundary is clean.
