package main

import "testing"

func TestPetHuntRemainUsed(t *testing.T) {
	const limit int64 = 10
	cases := []struct {
		name           string
		f1, f2         int64
		wantRemain     int64
		wantUsed       int64
		wantCanDraw    bool
	}{
		{"remaining full", 10, 0, 10, 0, true},
		{"remaining nine", 9, 0, 9, 1, true},
		{"remaining zero", 0, 10, 0, 10, false},
		{"lifetime in f1", 36, 3, 3, 7, true},
		{"lifetime in f2", 3, 36, 3, 7, true},
		{"unset both", 0, 0, 10, 0, true},
	}
	for _, c := range cases {
		remain, used, _ := petHuntRemainUsed(c.f1, c.f2, limit)
		if remain != c.wantRemain || used != c.wantUsed {
			t.Fatalf("%s: remain/used=%d/%d want %d/%d", c.name, remain, used, c.wantRemain, c.wantUsed)
		}
		if (remain > 0) != c.wantCanDraw {
			t.Fatalf("%s: canDraw=%v want %v", c.name, remain > 0, c.wantCanDraw)
		}
	}
}
