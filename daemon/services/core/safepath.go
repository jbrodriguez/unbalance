package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"unbalance/daemon/domain"
)

func cleanRoot(root string) (string, error) {
	if root == "" {
		return "", fmt.Errorf("empty root")
	}

	if !filepath.IsAbs(root) {
		return "", fmt.Errorf("root must be absolute: %s", root)
	}

	return filepath.Clean(root), nil
}

func cleanEntry(entry string) (string, error) {
	if entry == "" {
		return "", fmt.Errorf("empty entry")
	}

	if filepath.IsAbs(entry) {
		return "", fmt.Errorf("entry must be relative: %s", entry)
	}

	cleaned := filepath.Clean(entry)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("entry escapes root: %s", entry)
	}

	return cleaned, nil
}

func safeJoin(root, entry string) (string, string, error) {
	cleanedRoot, err := cleanRoot(root)
	if err != nil {
		return "", "", err
	}

	cleanedEntry, err := cleanEntry(entry)
	if err != nil {
		return "", "", err
	}

	joined := filepath.Join(cleanedRoot, cleanedEntry)
	if !pathIsInside(cleanedRoot, joined) {
		return "", "", fmt.Errorf("path escapes root: %s", joined)
	}

	return joined, cleanedEntry, nil
}

func pathIsInside(root, target string) bool {
	root = filepath.Clean(root)
	target = filepath.Clean(target)

	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}

	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func validateSymlinkBoundary(root, target string) error {
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("unable to resolve root %s: %w", root, err)
	}

	resolvedTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		return fmt.Errorf("unable to resolve path %s: %w", target, err)
	}

	if !pathIsInside(resolvedRoot, resolvedTarget) {
		return fmt.Errorf("path resolves outside root: %s", target)
	}

	return nil
}

type safeTransferPaths struct {
	SrcRoot string
	DstRoot string
	Entry   string
	SrcPath string
	DstPath string
}

func safeCommandPaths(command *domain.Command) (safeTransferPaths, error) {
	return safeCommandPathsWithin(command, nil)
}

func safeCommandPathsWithin(command *domain.Command, allowedRoots []string) (safeTransferPaths, error) {
	if command == nil {
		return safeTransferPaths{}, fmt.Errorf("missing command")
	}

	srcPath, entry, err := safeJoin(command.Src, command.Entry)
	if err != nil {
		return safeTransferPaths{}, fmt.Errorf("invalid source path: %w", err)
	}

	dstPath, dstEntry, err := safeJoin(command.Dst, command.Entry)
	if err != nil {
		return safeTransferPaths{}, fmt.Errorf("invalid destination path: %w", err)
	}

	if entry != dstEntry {
		return safeTransferPaths{}, fmt.Errorf("source and destination entries differ after cleaning")
	}

	srcRoot, _ := cleanRoot(command.Src)
	dstRoot, _ := cleanRoot(command.Dst)
	if srcRoot == dstRoot {
		return safeTransferPaths{}, fmt.Errorf("source and destination roots must differ")
	}

	if len(allowedRoots) > 0 {
		if !rootIsAllowed(srcRoot, allowedRoots) {
			return safeTransferPaths{}, fmt.Errorf("source root is not an allowed disk root: %s", srcRoot)
		}

		if !rootIsAllowed(dstRoot, allowedRoots) {
			return safeTransferPaths{}, fmt.Errorf("destination root is not an allowed disk root: %s", dstRoot)
		}
	}

	return safeTransferPaths{
		SrcRoot: srcRoot,
		DstRoot: dstRoot,
		Entry:   entry,
		SrcPath: srcPath,
		DstPath: dstPath,
	}, nil
}

func rootIsAllowed(root string, allowedRoots []string) bool {
	cleanedRoot, err := cleanRoot(root)
	if err != nil {
		return false
	}

	for _, allowed := range allowedRoots {
		cleanedAllowed, err := cleanRoot(allowed)
		if err != nil {
			continue
		}

		if cleanedRoot == cleanedAllowed {
			return true
		}
	}

	return false
}

func validateExistingPathBoundary(root, target string) error {
	if err := validateSymlinkBoundary(root, target); err != nil {
		return err
	}

	if _, err := os.Lstat(target); err != nil {
		return fmt.Errorf("unable to inspect path %s: %w", target, err)
	}

	return nil
}

func validateCreatablePathBoundary(root, target string) error {
	root, err := cleanRoot(root)
	if err != nil {
		return err
	}

	target = filepath.Clean(target)
	if !pathIsInside(root, target) {
		return fmt.Errorf("path escapes root: %s", target)
	}

	current := target
	for {
		_, statErr := os.Lstat(current)
		if statErr == nil {
			return validateSymlinkBoundary(root, current)
		}

		if !os.IsNotExist(statErr) {
			return fmt.Errorf("unable to inspect path %s: %w", current, statErr)
		}

		if current == root {
			return fmt.Errorf("root does not exist: %s", root)
		}

		current = filepath.Dir(current)
	}
}

func validateCommandForExecution(command *domain.Command, allowedRoots []string) (safeTransferPaths, error) {
	if len(allowedRoots) == 0 {
		return safeTransferPaths{}, fmt.Errorf("no allowed disk roots available")
	}

	paths, err := safeCommandPathsWithin(command, allowedRoots)
	if err != nil {
		return safeTransferPaths{}, err
	}

	if err := validateExistingPathBoundary(paths.SrcRoot, paths.SrcPath); err != nil {
		return paths, fmt.Errorf("invalid source path: %w", err)
	}

	if err := validateCreatablePathBoundary(paths.DstRoot, paths.DstPath); err != nil {
		return paths, fmt.Errorf("invalid destination path: %w", err)
	}

	return paths, nil
}

func (c *Core) allowedTransferRoots() []string {
	if c == nil || c.state == nil || c.state.Unraid == nil {
		return nil
	}

	roots := make([]string, 0, len(c.state.Unraid.Disks))
	for _, disk := range c.state.Unraid.Disks {
		if disk == nil || disk.Path == "" {
			continue
		}

		roots = append(roots, disk.Path)
	}

	return roots
}

func (c *Core) prepareCommandForExecution(command *domain.Command) (safeTransferPaths, error) {
	return c.transferExecutor().PrepareCommand(command, c.allowedTransferRoots())
}

func removeTransferredSource(command *domain.Command, pruneParents bool) (string, []string, error) {
	paths, err := safeCommandPaths(command)
	if err != nil {
		return "", nil, err
	}

	if _, err := os.Stat(paths.DstPath); err != nil {
		if os.IsNotExist(err) {
			return paths.SrcPath, nil, fmt.Errorf("destination is missing; refusing source deletion: %s", paths.DstPath)
		}

		return paths.SrcPath, nil, fmt.Errorf("unable to inspect destination %s: %w", paths.DstPath, err)
	}

	if err := validateSymlinkBoundary(paths.SrcRoot, paths.SrcPath); err != nil {
		return paths.SrcPath, nil, err
	}

	if err := validateSymlinkBoundary(paths.DstRoot, paths.DstPath); err != nil {
		return paths.SrcPath, nil, err
	}

	if err := os.RemoveAll(paths.SrcPath); err != nil {
		return paths.SrcPath, nil, err
	}

	if !pruneParents {
		return paths.SrcPath, nil, nil
	}

	pruned, err := pruneEmptyParents(paths.SrcRoot, paths.Entry)
	return paths.SrcPath, pruned, err
}

func pruneEmptyParents(srcRoot, entry string) ([]string, error) {
	root, err := cleanRoot(srcRoot)
	if err != nil {
		return nil, err
	}

	cleanedEntry, err := cleanEntry(entry)
	if err != nil {
		return nil, err
	}

	parentEntry := filepath.Dir(cleanedEntry)
	// Preserve historical behavior: skip user shares and top-level share children.
	if parentEntry == "." || !strings.Contains(parentEntry, string(filepath.Separator)) {
		return nil, nil
	}

	parent := filepath.Join(root, parentEntry)
	if !pathIsInside(root, parent) {
		return nil, fmt.Errorf("parent escapes root: %s", parent)
	}

	pruned := make([]string, 0)
	for pathIsInside(root, parent) && parent != root {
		rel, relErr := filepath.Rel(root, parent)
		if relErr != nil || rel == "." || !strings.Contains(rel, string(filepath.Separator)) {
			break
		}

		if err := validateSymlinkBoundary(root, parent); err != nil {
			return pruned, err
		}

		err := os.Remove(parent)
		if err == nil {
			pruned = append(pruned, parent)
			parent = filepath.Dir(parent)
			continue
		}

		if os.IsNotExist(err) {
			parent = filepath.Dir(parent)
			continue
		}

		if os.IsPermission(err) {
			return pruned, err
		}

		// Non-empty directories stop pruning. This is expected and safe.
		return pruned, nil
	}

	return pruned, nil
}
