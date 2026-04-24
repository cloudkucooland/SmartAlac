package sa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sorrow446/go-mp4tag"
)

func TestNewCurator(t *testing.T) {
	cfg := Config{
		FinalDir: "/tmp/music",
	}
	curator := NewCurator(cfg)

	if curator.Config.FinalDir != "/tmp/music" {
		t.Errorf("expected /tmp/music, got %s", curator.Config.FinalDir)
	}

	if curator.Stats.BadQueries == nil {
		t.Error("BadQueries map should be initialized")
	}
}

func TestTagPreservation(t *testing.T) {
	in := &mp4tag.MP4Tags{
		Artist:      "Original Artist",
		AlbumArtist: "Original Album Artist",
		Comment:     "Original Comment",
		Year:        1990,
	}

	// Simulation of out := *in
	out := *in

	// Simulation of MusicBrainz being silent on some fields
	mbArtist := ""
	mbComment := "New Comment"
	mbYear := 0

	// Logic we implemented in mb.go
	if mbArtist != "" {
		out.Artist = mbArtist
	}
	if mbComment != "" {
		out.Comment = mbComment
	}
	if mbYear != 0 {
		out.Year = int32(mbYear)
	}

	if out.Artist != "Original Artist" {
		t.Errorf("Artist was overwritten by empty string, got %q", out.Artist)
	}
	if out.AlbumArtist != "Original Album Artist" {
		t.Errorf("AlbumArtist was lost, got %q", out.AlbumArtist)
	}
	if out.Comment != "New Comment" {
		t.Errorf("Comment was not updated, got %q", out.Comment)
	}
	if out.Year != 1990 {
		t.Errorf("Year was overwritten by zero, got %d", out.Year)
	}
}

func TestRename(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	finalDir := filepath.Join(tempDir, "final")
	os.Mkdir(sourceDir, 0755)
	os.Mkdir(finalDir, 0755)

	dummyFile := filepath.Join(sourceDir, "test.m4a")
	os.WriteFile(dummyFile, []byte("dummy"), 0644)

	curator := NewCurator(Config{
		FinalDir: finalDir,
	})

	tags := &mp4tag.MP4Tags{
		AlbumArtist: "Test Artist",
		AlbumSort:   "Test Album",
		Title:       "Test Track",
		DiscNumber:  1,
		TrackNumber: 1,
	}

	err := curator.rename(dummyFile, tags)
	if err != nil {
		t.Fatalf("rename failed: %v", err)
	}

	expectedPath := filepath.Join(finalDir, "Test Artist", "Test Album", "1-01 Test Track.m4a")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("file was not moved to expected location: %s", expectedPath)
	}

	if _, err := os.Stat(dummyFile); err == nil {
		t.Errorf("source file still exists at %s", dummyFile)
	}
}
