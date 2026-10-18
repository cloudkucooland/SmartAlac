package bme

import (
	"bytes"
	"encoding/binary"
	"github.com/cloudkucooland/SmartAlac/pkg/cdio"
	"testing"
)

func TestWriteWavHeader(t *testing.T) {
	var buf bytes.Buffer
	size := uint32(1000)
	write_wav_header(&buf, size)

	data := buf.Bytes()
	if len(data) != 44 {
		t.Errorf("expected header size 44, got %d", len(data))
	}

	if string(data[0:4]) != "RIFF" {
		t.Errorf("expected RIFF, got %s", string(data[0:4]))
	}

	riffSize := binary.LittleEndian.Uint32(data[4:8])
	if riffSize != size+44-8 {
		t.Errorf("expected riff size %d, got %d", size+44-8, riffSize)
	}

	if string(data[8:12]) != "WAVE" {
		t.Errorf("expected WAVE, got %s", string(data[8:12]))
	}

	if string(data[12:16]) != "fmt " {
		t.Errorf("expected 'fmt ', got %s", string(data[12:16]))
	}

	dataSize := binary.LittleEndian.Uint32(data[40:44])
	if dataSize != size {
		t.Errorf("expected data size %d, got %d", size, dataSize)
	}
}

func TestGetMbdiscid(t *testing.T) {
	// Mock the cdio functions used by get_mbdiscid
	old_get_first_track_num := cdio.GetFirstTrackNum
	old_get_num_tracks := cdio.GetNumTracks
	old_get_track_lba := cdio.GetTrackLba
	defer func() {
		cdio.GetFirstTrackNum = old_get_first_track_num
		cdio.GetNumTracks = old_get_num_tracks
		cdio.GetTrackLba = old_get_track_lba
	}()

	// Example from MusicBrainz documentation or known values
	// Let's use a simple 3-track CD example
	cdio.GetFirstTrackNum = func(d cdio.Device) cdio.Track { return 1 }
	cdio.GetNumTracks = func(d cdio.Device) cdio.Track { return 3 }

	lbas := map[cdio.Track]int{
		1:                 150,
		2:                 15000,
		3:                 30000,
		cdio.LeadoutTrack: 45000,
	}

	cdio.GetTrackLba = func(d cdio.Device, track cdio.Track) int {
		return lbas[track]
	}

	discid := get_mbdiscid(nil)
	if discid == "" {
		t.Fatal("expected a discid, got empty string")
	}

	t.Logf("Generated DiscID: %s", discid)

	for _, char := range discid {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '.' && char != '_' && char != '-' {
			t.Errorf("discid contains invalid character: %c", char)
		}
	}
}
