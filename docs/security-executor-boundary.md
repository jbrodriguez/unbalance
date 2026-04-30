# Executor boundary details

Branch: `security/root-runtime-review`

This expands the proposed phase 1 work: create an internal filesystem executor boundary before splitting into separate web/worker processes.

## 1. `filesystem.Executor` boundary

The first refactor should make privileged filesystem actions explicit behind an interface.

Example shape:

```go
type Executor interface {
    StorageStatus(ctx context.Context) (*domain.Unraid, error)
    ListTree(ctx context.Context, req ListTreeRequest) (domain.Branch, error)
    Locate(ctx context.Context, req LocateRequest) ([]string, error)
    ScanItems(ctx context.Context, req ScanItemsRequest) (*ScanItemsResult, error)
    ExecuteTransfer(ctx context.Context, req TransferRequest, events chan<- TransferEvent) (*TransferResult, error)
    ValidateTransfer(ctx context.Context, req TransferRequest, events chan<- TransferEvent) (*TransferResult, error)
    RemoveSource(ctx context.Context, req RemoveSourceRequest) error
    Stop(ctx context.Context, operationID string) error
}
```

Initially this can be an in-process implementation using the existing `os`, `find`, `du`, and `rsync` logic. Later the same interface can be backed by a Unix-socket worker client.

The important win is architectural: core/web code stops directly shelling out or deleting paths.

## 2. Move semantics: copy first, then guarded delete

A move plan should remain logically the same:

```text
1. rsync/copy source entry to destination
2. verify copy completed acceptably
3. delete only the exact source entry that was copied
4. optionally prune empty parents, bounded by safety rules
```

What changes is where the delete authority lives and how it is validated.

### Current shape

Current code effectively does:

```text
rsync src/entry dst/
if dst/entry exists:
    shell: rm -rf src/entry
    maybe shell: find parent -type d -empty -prune -exec rm -rf {}
```

### Safer shape

The executor/worker receives a structured command:

```go
type TransferCommand struct {
    ID       string
    SrcRoot  string // e.g. /mnt/disk1
    DstRoot  string // e.g. /mnt/disk2
    Entry    string // e.g. media/movie/file.mkv, always relative
    Mode     TransferMode // Copy or Move
    RsyncArgs []string
}
```

For `ModeMove`, the executor runs:

```text
safeSrc = validateJoin(SrcRoot, Entry)
safeDst = validateJoin(DstRoot, Entry)
run rsync safeSrc -> DstRoot
if rsync success or acceptable partial status:
    verify safeDst exists
    remove safeSrc only if it still validates and matches the completed command
    prune empty parents upward, stopping before SrcRoot/share boundary
```

So deletion is not a separate browser-triggered “delete this path” ability. It is a postcondition of a specific completed transfer command.

### Important delete rules

For move deletion:

- delete only after the matching transfer completes;
- delete only `SrcRoot + Entry`, never arbitrary client paths;
- revalidate immediately before deletion;
- refuse `Entry` values that are absolute, empty, `.`, `..`, or escape via `..`;
- refuse source roots outside discovered array disks;
- refuse destination roots outside discovered array disks;
- refuse same source/destination root;
- do not follow symlink escapes for destructive deletion;
- if verification is uncertain, skip delete and flag the command.

This preserves current move behavior while tightening the authority boundary.

## 3. Replacing shell deletion with Go deletion

Instead of building shell strings like:

```bash
rm -rf "/mnt/disk1/share/item"
find "/mnt/disk1/share/parent" -type d -empty -prune -exec rm -rf {} \;
```

use guarded Go functions.

Example sketch:

```go
func RemoveTransferredSource(cmd TransferCommand, roots AllowedRoots) error {
    src, err := ValidateEntryPath(roots, cmd.SrcRoot, cmd.Entry)
    if err != nil {
        return err
    }

    dst, err := ValidateEntryPath(roots, cmd.DstRoot, cmd.Entry)
    if err != nil {
        return err
    }

    if !Exists(dst) {
        return fmt.Errorf("destination is missing; refusing source deletion")
    }

    if IsUnsafeRemovalTarget(src, cmd.SrcRoot) {
        return fmt.Errorf("unsafe removal target")
    }

    return os.RemoveAll(src)
}
```

Parent pruning should be intentionally conservative:

```go
func PruneEmptyParents(srcRoot, entry string) {
    parent := filepath.Dir(filepath.Join(srcRoot, entry))

    for isBelowBoundary(parent, srcRoot) {
        err := os.Remove(parent) // succeeds only if empty
        if err != nil {
            return
        }
        parent = filepath.Dir(parent)
    }
}
```

Use `os.Remove`, not `os.RemoveAll`, for parent pruning. That way non-empty directories are naturally preserved.

The result is safer than shelling out because:

- no shell parser;
- no command injection class;
- easier boundary checks;
- easier unit tests;
- clearer audit logs;
- no accidental broad `rm -rf` string construction.

## 4. Server-side IDs instead of executable browser payloads

The browser should not be trusted to send executable operations.

### Current risk

Some flows accept an entire `Operation` or `Command` object from the browser/websocket. Those objects contain fields that can affect filesystem behavior:

- `Src`
- `Dst`
- `Entry`
- `RsyncArgs`
- `OpKind`

Even with auth, the browser is the least trusted part of the system. It should display and confirm plans, not define executable paths.

### Safer model

The server creates and stores immutable plans.

```go
type StoredPlan struct {
    ID        string
    Kind      OperationKind
    CreatedAt time.Time
    Commands  []TransferCommand
    Summary   PlanSummary
    State     PlanState
}
```

Flow:

```text
Browser -> web: plan request { source, targets, selected folders }
Web/core: normalizes + validates request
Web/core: creates StoredPlan with server-generated commands
Web -> browser: display plan summary + planID
Browser -> web: execute { planID }
Web: loads StoredPlan(planID)
Web -> executor: ExecuteTransfer(storedPlan.Commands)
```

For individual source removal after a flagged or manual flow:

```text
Browser -> web: removeSource { planID, commandID }
Web: finds exact stored command
Web: verifies command status allows removal
Web -> executor: RemoveSource(storedCommand)
```

For replay:

```text
Browser -> web: replay { historyOperationID }
Web: loads server-side history record
Web: creates new StoredPlan from that trusted record after revalidation
Web -> browser: shows revalidated plan
Browser -> web: execute { newPlanID }
```

This means a malicious browser cannot invent:

```json
{ "src": "/", "entry": "boot", "dst": "/mnt/disk1" }
```

because the server never executes browser-provided `src/dst/entry`. It only executes commands it created and stored.

## Practical first PR scope

A good first PR could be small and still valuable:

1. Add `daemon/services/filesystem` or `daemon/services/executor` package.
2. Add path validation helpers + tests.
3. Move source deletion/pruning behind executor functions.
4. Replace shell `rm -rf` and shell parent pruning with guarded Go code.
5. Leave process model unchanged.

Then a second PR can tackle server-side `planID` execution.
