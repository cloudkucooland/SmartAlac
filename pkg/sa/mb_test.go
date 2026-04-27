package sa

import (
	"github.com/Sorrow446/go-mp4tag"
	"testing"
)

func TestInferTrackInfo(t *testing.T) {
	cases := []struct {
		name      string
		wantDisc  int16
		wantTrack int16
	}{
		{"01 National Brotherhood Week.wav.m4a", 0, 1},
		{"1-13 New Math.m4a", 1, 13},
		{"7 Pollution.m4a", 0, 7},
		{"1-7 Something.m4a", 1, 7},
		{"08 So Long, Mom (A Song for World War III).wav.m4a", 0, 8},
		{"NoNumber.m4a", 0, 0},
		{"  09  Whitespace.m4a", 0, 9},
		{"2-01 CD2 Track 1.m4a", 2, 1},
		{"invalid-format.m4a", 0, 0},
		{"10-25 Huge Numbers.m4a", 10, 25},
		{"-5 Negative.m4a", 0, 5},
	}

	for _, c := range cases {
		gotDisc, gotTrack := inferTrackInfo(c.name)
		if gotDisc != c.wantDisc || gotTrack != c.wantTrack {
			t.Errorf("inferTrackInfo(%q) == (%d, %d), want (%d, %d)", c.name, gotDisc, gotTrack, c.wantDisc, c.wantTrack)
		}
	}
}

func TestTagCleanup(t *testing.T) {
	c := &Curator{}

	t.Run("cleanupOtherCustom", func(t *testing.T) {
		tags := &mp4tag.MP4Tags{
			Custom: map[string]string{
				"ENGINEER": "Dave Cook",
			},
			OtherCustom: map[string][]string{
				"ENGINEER":      {"Dennice Brown", "Dave Cook", "Frank Filipetti"},
				"CATALOGNUMBER": {"9 60815-2"},
			},
		}

		c.cleanupOtherCustom(tags)

		if _, ok := tags.OtherCustom["ENGINEER"]; ok {
			t.Error("ENGINEER should have been removed from OtherCustom because it exists in Custom")
		}
		if _, ok := tags.OtherCustom["CATALOGNUMBER"]; !ok {
			t.Error("CATALOGNUMBER should have been preserved in OtherCustom")
		}
	})

	t.Run("uniquifyOtherCustom", func(t *testing.T) {
		tags := &mp4tag.MP4Tags{
			Custom: map[string]string{
				"LABEL": "Elektra",
			},
			OtherCustom: map[string][]string{
				"MIXER":         {"Frank Wolf", "Frank Wolf", "Frank Wolf"},
				"CATALOGNUMBER": {"9 60738-1", "9 60738-1"},
				"LABEL":         {"Elektra", "Warner"},
			},
		}

		c.uniquifyOtherCustom(tags)

		if len(tags.OtherCustom["MIXER"]) != 1 || tags.OtherCustom["MIXER"][0] != "Frank Wolf" {
			t.Errorf("MIXER should be uniquely Frank Wolf, got %v", tags.OtherCustom["MIXER"])
		}
		if len(tags.OtherCustom["CATALOGNUMBER"]) != 1 || tags.OtherCustom["CATALOGNUMBER"][0] != "9 60738-1" {
			t.Errorf("CATALOGNUMBER should be uniquely 9 60738-1, got %v", tags.OtherCustom["CATALOGNUMBER"])
		}
		// LABEL in OtherCustom should have "Elektra" removed because it is in Custom, and "Warner" should remain.
		if len(tags.OtherCustom["LABEL"]) != 1 || tags.OtherCustom["LABEL"][0] != "Warner" {
			t.Errorf("LABEL in OtherCustom should have Warner left, got %v", tags.OtherCustom["LABEL"])
		}
	})
}
