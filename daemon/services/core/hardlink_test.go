package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"unbalance/daemon/domain"
)

// TestEnvironment manages temporary test directories and files
type TestEnvironment struct {
	TempDir string
	Files   map[string]string // filename -> content
}

// NewTestEnvironment creates a new test environment with temporary directory
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	tempDir, err := os.MkdirTemp("", "unbalance_hardlink_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	return &TestEnvironment{
		TempDir: tempDir,
		Files:   make(map[string]string),
	}
}

// Cleanup removes the temporary test environment
func (te *TestEnvironment) Cleanup() {
	os.RemoveAll(te.TempDir)
}

// CreateFile creates a file with given content in the test environment
func (te *TestEnvironment) CreateFile(filename, content string) (string, error) {
	fullPath := filepath.Join(te.TempDir, filename)
	
	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	
	// Write file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", err
	}
	
	te.Files[filename] = content
	return fullPath, nil
}

// CreateHardlink creates a hardlink from source to destination
func (te *TestEnvironment) CreateHardlink(src, dst string) error {
	srcPath := filepath.Join(te.TempDir, src)
	dstPath := filepath.Join(te.TempDir, dst)
	
	// Create directory for destination if needed
	dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.Link(srcPath, dstPath)
}

// CreateDirectory creates a directory in the test environment
func (te *TestEnvironment) CreateDirectory(dirname string) error {
	fullPath := filepath.Join(te.TempDir, dirname)
	return os.MkdirAll(fullPath, 0755)
}

// GetPath returns the full path for a relative filename
func (te *TestEnvironment) GetPath(filename string) string {
	return filepath.Join(te.TempDir, filename)
}

func TestDetectHardlinks_SingleFile(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create a regular file
	filePath, err := env.CreateFile("regular.txt", "This is a regular file")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Test detection
	info, err := detectHardlinks(filePath)
	if err != nil {
		t.Fatalf("detectHardlinks failed: %v", err)
	}

	if info.HasHardlinks {
		t.Errorf("Regular file should not have hardlinks")
	}
	if info.InodeCount != 1 {
		t.Errorf("Regular file should have link count of 1, got %d", info.InodeCount)
	}
	if info.SharedSize != info.ExpandedSize {
		t.Errorf("Regular file shared size (%d) should equal expanded size (%d)", 
			info.SharedSize, info.ExpandedSize)
	}
}

func TestDetectHardlinks_WithHardlinks(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create original file
	content := "This file has hardlinks"
	originalPath, err := env.CreateFile("original.txt", content)
	if err != nil {
		t.Fatalf("Failed to create original file: %v", err)
	}

	// Create hardlinks
	if err := env.CreateHardlink("original.txt", "link1.txt"); err != nil {
		t.Fatalf("Failed to create hardlink 1: %v", err)
	}
	if err := env.CreateHardlink("original.txt", "link2.txt"); err != nil {
		t.Fatalf("Failed to create hardlink 2: %v", err)
	}

	// Test detection on original file
	info, err := detectHardlinks(originalPath)
	if err != nil {
		t.Fatalf("detectHardlinks failed: %v", err)
	}

	if !info.HasHardlinks {
		t.Errorf("File should have hardlinks")
	}
	if info.InodeCount != 3 {
		t.Errorf("File should have link count of 3, got %d", info.InodeCount)
	}
	
	expectedSize := uint64(len(content))
	if info.SharedSize != expectedSize {
		t.Errorf("Shared size should be %d, got %d", expectedSize, info.SharedSize)
	}
	if info.ExpandedSize != expectedSize*3 {
		t.Errorf("Expanded size should be %d, got %d", expectedSize*3, info.ExpandedSize)
	}

	// Test detection on hardlink
	linkPath := env.GetPath("link1.txt")
	linkInfo, err := detectHardlinks(linkPath)
	if err != nil {
		t.Fatalf("detectHardlinks on hardlink failed: %v", err)
	}

	// Should detect same hardlink info
	if linkInfo.Inode != info.Inode {
		t.Errorf("Hardlink should have same inode as original")
	}
	if linkInfo.InodeCount != info.InodeCount {
		t.Errorf("Hardlink should have same link count as original")
	}
}

func TestDetectHardlinks_Directory(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create directory structure
	env.CreateDirectory("testdir")
	
	// Create files with mixed hardlink scenarios
	env.CreateFile("testdir/regular.txt", "regular file")
	env.CreateFile("testdir/hardlinked.txt", "hardlinked file content")
	env.CreateHardlink("testdir/hardlinked.txt", "testdir/hardlink1.txt")
	env.CreateHardlink("testdir/hardlinked.txt", "testdir/hardlink2.txt")
	
	// Create another set of hardlinks
	env.CreateFile("testdir/another.txt", "another hardlinked file")
	env.CreateHardlink("testdir/another.txt", "testdir/another_link.txt")

	// Test directory detection
	dirPath := env.GetPath("testdir")
	info, err := detectHardlinks(dirPath)
	if err != nil {
		t.Fatalf("detectHardlinks on directory failed: %v", err)
	}

	if !info.HasHardlinks {
		t.Errorf("Directory should contain hardlinks")
	}

	// Verify total sizes account for all hardlinks
	// We should have:
	// - 1 regular file (no hardlinks)
	// - 1 hardlinked file with 3 links total
	// - 1 another file with 2 links total
	expectedHardlinkCount := 5 // 3 + 2 total hardlinked files
	if info.InodeCount != expectedHardlinkCount {
		t.Errorf("Expected %d total hardlinks, got %d", expectedHardlinkCount, info.InodeCount)
	}

	// Verify expanded size is larger than shared size
	if info.ExpandedSize <= info.SharedSize {
		t.Errorf("Expanded size (%d) should be larger than shared size (%d)", 
			info.ExpandedSize, info.SharedSize)
	}
}

func TestGetEffectiveSize(t *testing.T) {
	// Create test item with hardlink info
	item := &domain.Item{
		Name: "test.txt",
		Size: 1000, // This will be updated by hardlink detection
		HardlinkInfo: &domain.HardlinkInfo{
			HasHardlinks: true,
			SharedSize:   1000,
			ExpandedSize: 3000, // 3 hardlinks
			InodeCount:   3,
		},
	}

	// Test with hardlinks preserved
	effectiveSize := getEffectiveSize(item, true)
	if effectiveSize != 1000 {
		t.Errorf("With hardlinks preserved, effective size should be 1000, got %d", effectiveSize)
	}

	// Test with hardlinks not preserved (conservative)
	effectiveSize = getEffectiveSize(item, false)
	if effectiveSize != 3000 {
		t.Errorf("Without hardlinks preserved, effective size should be 3000, got %d", effectiveSize)
	}

	// Test with item that has no hardlinks
	regularItem := &domain.Item{
		Name: "regular.txt",
		Size: 500,
		HardlinkInfo: &domain.HardlinkInfo{
			HasHardlinks: false,
			SharedSize:   500,
			ExpandedSize: 500,
			InodeCount:   1,
		},
	}

	effectiveSize = getEffectiveSize(regularItem, true)
	if effectiveSize != 500 {
		t.Errorf("Regular file effective size should be 500, got %d", effectiveSize)
	}

	effectiveSize = getEffectiveSize(regularItem, false)
	if effectiveSize != 500 {
		t.Errorf("Regular file effective size should be 500, got %d", effectiveSize)
	}
}

func TestUpdateItemWithHardlinkInfo(t *testing.T) {
	// Test item with hardlinks
	item := &domain.Item{
		Name: "hardlinked.txt",
		Size: 1000,
		HardlinkInfo: &domain.HardlinkInfo{
			HasHardlinks: true,
			SharedSize:   1000,
			ExpandedSize: 4000, // 4 hardlinks
			InodeCount:   4,
		},
	}

	// Test with hardlinks preserved
	updateItemWithHardlinkInfo(item, true)
	if item.Size != 1000 {
		t.Errorf("With hardlinks preserved, item size should be 1000, got %d", item.Size)
	}

	// Reset item size and test without hardlinks preserved
	item.Size = 1000
	updateItemWithHardlinkInfo(item, false)
	if item.Size != 4000 {
		t.Errorf("Without hardlinks preserved, item size should be 4000, got %d", item.Size)
	}
}

func TestHardlinkDetection_EdgeCases(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Test with empty file
	emptyPath, err := env.CreateFile("empty.txt", "")
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	info, err := detectHardlinks(emptyPath)
	if err != nil {
		t.Fatalf("detectHardlinks on empty file failed: %v", err)
	}

	if info.SharedSize != 0 {
		t.Errorf("Empty file should have size 0, got %d", info.SharedSize)
	}

	// Test with non-existent file
	_, err = detectHardlinks("/nonexistent/path")
	if err == nil {
		t.Errorf("detectHardlinks should fail on non-existent file")
	}

	// Test with large file
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	
	largePath, err := env.CreateFile("large.txt", string(largeContent))
	if err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	// Create hardlinks for large file
	env.CreateHardlink("large.txt", "large_link1.txt")
	env.CreateHardlink("large.txt", "large_link2.txt")

	info, err = detectHardlinks(largePath)
	if err != nil {
		t.Fatalf("detectHardlinks on large file failed: %v", err)
	}

	if !info.HasHardlinks {
		t.Errorf("Large file should have hardlinks")
	}
	if info.SharedSize != uint64(len(largeContent)) {
		t.Errorf("Large file shared size mismatch")
	}
	if info.ExpandedSize != uint64(len(largeContent))*3 {
		t.Errorf("Large file expanded size mismatch")
	}
}

// Benchmark tests for performance validation
func BenchmarkDetectHardlinks_SingleFile(b *testing.B) {
	env := NewTestEnvironment(&testing.T{})
	defer env.Cleanup()

	filePath, _ := env.CreateFile("bench.txt", "benchmark file content")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectHardlinks(filePath)
	}
}

func BenchmarkDetectHardlinks_ManyHardlinks(b *testing.B) {
	env := NewTestEnvironment(&testing.T{})
	defer env.Cleanup()

	// Create original file
	env.CreateFile("original.txt", "content for many hardlinks")

	// Create many hardlinks
	for i := 0; i < 100; i++ {
		env.CreateHardlink("original.txt", fmt.Sprintf("link_%d.txt", i))
	}

	filePath := env.GetPath("original.txt")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectHardlinks(filePath)
	}
}

func BenchmarkDetectHardlinks_Directory(b *testing.B) {
	env := NewTestEnvironment(&testing.T{})
	defer env.Cleanup()

	// Create directory with mixed files
	env.CreateDirectory("benchdir")
	
	// Create various file scenarios
	for i := 0; i < 50; i++ {
		content := fmt.Sprintf("file content %d", i)
		filename := fmt.Sprintf("benchdir/file_%d.txt", i)
		env.CreateFile(filename, content)
		
		// Make some files hardlinked
		if i%3 == 0 {
			env.CreateHardlink(filename, fmt.Sprintf("benchdir/link_%d.txt", i))
		}
	}

	dirPath := env.GetPath("benchdir")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectHardlinks(dirPath)
	}
}