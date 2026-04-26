package sa

import "testing"

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
		{"-5 Negative.m4a", 0, 5}, // Splits into "" and "5", strconv.Atoi("") fails for disc, succeeds for track
	}

	for _, c := range cases {
		gotDisc, gotTrack := inferTrackInfo(c.name)
		if gotDisc != c.wantDisc || gotTrack != c.wantTrack {
			t.Errorf("inferTrackInfo(%q) == (%d, %d), want (%d, %d)", c.name, gotDisc, gotTrack, c.wantDisc, c.wantTrack)
		}
	}
}
