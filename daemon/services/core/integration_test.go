package core

import (
	"testing"

	"unbalance/daemon/domain"
)

// IntegrationTestEnvironment provides a more comprehensive test environment
// that simulates the full unbalanced workflow
type IntegrationTestEnvironment struct {
	*TestEnvironment
	Config *domain.Config
}

// NewIntegrationTestEnvironment creates a test environment with configuration
func NewIntegrationTestEnvironment(t *testing.T, preserveHardlinks bool) *IntegrationTestEnvironment {
	env := NewTestEnvironment(t)
	
	config := &domain.Config{
		PreserveHardlinks: preserveHardlinks,
		DryRun:           true,
		ReservedAmount:   1024 * 1024 * 1024, // 1GB
		ReservedUnit:     "bytes",
		RsyncArgs:        []string{"-X"},
	}
	
	return &IntegrationTestEnvironment{
		TestEnvironment: env,
		Config:         config,
	}
}

// CreateComplexHardlinkStructure creates a realistic directory structure with hardlinks
func (env *IntegrationTestEnvironment) CreateComplexHardlinkStructure() error {
	// Create directory structure
	dirs := []string{
		"media/movies",
		"media/tv",
		"documents/projects",
		"documents/archives",
		"backup/daily",
		"backup/weekly",
	}
	
	for _, dir := range dirs {
		if err := env.CreateDirectory(dir); err != nil {
			return err
		}
	}
	
	// Create regular files
	regularFiles := map[string]string{
		"media/movies/movie1.mkv":      string(make([]byte, 2*1024*1024*1024)), // 2GB
		"media/movies/movie2.mkv":      string(make([]byte, 1*1024*1024*1024)), // 1GB
		"media/tv/episode1.mkv":        string(make([]byte, 500*1024*1024)),    // 500MB
		"documents/projects/readme.txt": "Project documentation",
		"documents/archives/old.tar":   string(make([]byte, 100*1024*1024)),    // 100MB
	}
	
	for filename, content := range regularFiles {
		if _, err := env.CreateFile(filename, content); err != nil {
			return err
		}
	}
	
	// Create hardlinked backup files (common scenario)
	backupFiles := map[string]string{
		"backup/original/config.json":   `{"setting": "value"}`,
		"backup/original/database.db":   string(make([]byte, 50*1024*1024)), // 50MB
		"backup/original/logs.txt":      "application logs",
	}
	
	// Create original backup files
	env.CreateDirectory("backup/original")
	for filename, content := range backupFiles {
		if _, err := env.CreateFile(filename, content); err != nil {
			return err
		}
	}
	
	// Create hardlinks in daily backup
	hardlinks := map[string]string{
		"backup/original/config.json": "backup/daily/config.json",
		"backup/original/database.db": "backup/daily/database.db",
		"backup/original/logs.txt":    "backup/daily/logs.txt",
	}
	
	for src, dst := range hardlinks {
		if err := env.CreateHardlink(src, dst); err != nil {
			return err
		}
	}
	
	// Create additional hardlinks in weekly backup
	weeklyHardlinks := map[string]string{
		"backup/original/config.json": "backup/weekly/config.json",
		"backup/original/database.db": "backup/weekly/database.db",
	}
	
	for src, dst := range weeklyHardlinks {
		if err := env.CreateHardlink(src, dst); err != nil {
			return err
		}
	}
	
	// Create a complex hardlink scenario - same file in multiple locations
	env.CreateFile("shared/important.txt", "very important data")
	sharedLinks := []string{
		"documents/projects/important.txt",
		"backup/daily/important.txt",
		"backup/weekly/important.txt",
		"media/important.txt",
	}
	
	for _, link := range sharedLinks {
		if err := env.CreateHardlink("shared/important.txt", link); err != nil {
			return err
		}
	}
	
	return nil
}

func TestDirectHardlinkDetection(t *testing.T) {
	env := NewIntegrationTestEnvironment(t, true) // Enable hardlink preservation
	defer env.Cleanup()
	
	// Create simple test structure
	env.CreateDirectory("testdir")
	env.CreateFile("testdir/file1.txt", "content1")
	env.CreateFile("testdir/file2.txt", "content2")
	env.CreateHardlink("testdir/file1.txt", "testdir/file1_link.txt")
	
	// Test directory hardlink detection
	dirPath := env.GetPath("testdir")
	info, err := detectHardlinks(dirPath)
	if err != nil {
		t.Fatalf("detectHardlinks failed: %v", err)
	}
	
	if !info.HasHardlinks {
		t.Error("Directory should contain hardlinks")
	}
	
	if info.InodeCount != 2 {
		t.Errorf("Expected 2 hardlinked files, got %d", info.InodeCount)
	}
	
	if info.ExpandedSize <= info.SharedSize {
		t.Errorf("Expanded size (%d) should be larger than shared size (%d)", 
			info.ExpandedSize, info.SharedSize)
	}
}

func TestAnalyzeHardlinks_Direct(t *testing.T) {
	// Create test items directly without relying on getItems
	items := []*domain.Item{
		{
			Name: "regular.txt",
			Size: 1000,
			HardlinkInfo: &domain.HardlinkInfo{
				HasHardlinks: false,
				SharedSize:   1000,
				ExpandedSize: 1000,
				InodeCount:   1,
			},
		},
		{
			Name: "hardlinked1.txt",
			Size: 2000,
			HardlinkInfo: &domain.HardlinkInfo{
				HasHardlinks: true,
				SharedSize:   2000,
				ExpandedSize: 6000, // 3 hardlinks
				InodeCount:   3,
			},
		},
		{
			Name: "hardlinked2.txt",
			Size: 1500,
			HardlinkInfo: &domain.HardlinkInfo{
				HasHardlinks: true,
				SharedSize:   1500,
				ExpandedSize: 3000, // 2 hardlinks
				InodeCount:   2,
			},
		},
	}
	
	// Analyze hardlinks
	hardlinkedCount, sharedSize, expandedSize := analyzeHardlinks(items)
	
	if hardlinkedCount != 2 { // 2 items that have hardlinks
		t.Errorf("Expected 2 items with hardlinks, got %d", hardlinkedCount)
	}
	
	expectedSharedSize := uint64(2000 + 1500) // Only hardlinked files
	if sharedSize != expectedSharedSize {
		t.Errorf("Expected shared size %d, got %d", expectedSharedSize, sharedSize)
	}
	
	expectedExpandedSize := uint64(6000 + 3000) // Expanded hardlinked files
	if expandedSize != expectedExpandedSize {
		t.Errorf("Expected expanded size %d, got %d", expectedExpandedSize, expandedSize)
	}
	
	if expandedSize <= sharedSize {
		t.Errorf("Expected expanded size (%d) to be larger than shared size (%d)", 
			expandedSize, sharedSize)
	}
}

func TestUpdateItemWithHardlinkInfo_Integration(t *testing.T) {
	env := NewIntegrationTestEnvironment(t, false) // Start with hardlinks disabled
	defer env.Cleanup()
	
	// Create a file with hardlinks
	content := "test file content for hardlink testing"
	env.CreateFile("original.txt", content)
	env.CreateHardlink("original.txt", "link1.txt")
	env.CreateHardlink("original.txt", "link2.txt")
	
	// Detect hardlinks
	originalPath := env.GetPath("original.txt")
	hardlinkInfo, err := detectHardlinks(originalPath)
	if err != nil {
		t.Fatalf("Failed to detect hardlinks: %v", err)
	}
	
	// Create item
	item := &domain.Item{
		Name:         "original.txt",
		Size:         uint64(len(content)),
		Path:         "original.txt",
		HardlinkInfo: hardlinkInfo,
	}
	
	// Test with hardlinks not preserved (conservative)
	updateItemWithHardlinkInfo(item, false)
	expectedExpandedSize := uint64(len(content)) * 3 // 3 hardlinks
	if item.Size != expectedExpandedSize {
		t.Errorf("Expected conservative size %d, got %d", expectedExpandedSize, item.Size)
	}
	
	// Test with hardlinks preserved
	item.Size = uint64(len(content)) // Reset
	updateItemWithHardlinkInfo(item, true)
	expectedSharedSize := uint64(len(content))
	if item.Size != expectedSharedSize {
		t.Errorf("Expected shared size %d, got %d", expectedSharedSize, item.Size)
	}
}

func TestRealWorldScenario_BackupDirectory(t *testing.T) {
	env := NewIntegrationTestEnvironment(t, true)
	defer env.Cleanup()
	
	// Create a simple backup scenario
	env.CreateFile("backup1/file.txt", "content")
	env.CreateHardlink("backup1/file.txt", "backup2/file.txt")
	env.CreateHardlink("backup1/file.txt", "backup3/file.txt")
	
	// Test hardlink detection on the directory structure
	info1, err := detectHardlinks(env.GetPath("backup1"))
	if err != nil {
		t.Fatalf("Failed to detect hardlinks in backup1: %v", err)
	}
	
	if !info1.HasHardlinks {
		t.Error("backup1 should contain hardlinks")
	}
	
	// Verify space savings calculation
	if info1.ExpandedSize <= info1.SharedSize {
		t.Errorf("Expanded size (%d) should be larger than shared size (%d)", 
			info1.ExpandedSize, info1.SharedSize)
	}
	
	spaceSavings := info1.ExpandedSize - info1.SharedSize
	t.Logf("Backup scenario space savings: %d bytes", spaceSavings)
	
	if spaceSavings == 0 {
		t.Error("Expected space savings from hardlink preservation")
	}
}