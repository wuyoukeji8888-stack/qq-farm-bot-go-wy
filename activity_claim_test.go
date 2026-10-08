package main

import "testing"

func TestActEnsureGroupRootsFromChild(t *testing.T) {
	items := actEnsureGroupRoots([]*ActivityInfo{
		{ID: actWishSignID, Title: "秋祈良愿", StartTime: 1, EndTime: 9},
		{ID: actJoyShareRoot + 1, Title: "快乐不独享", StartTime: 2, EndTime: 8},
	})
	have := map[int64]string{}
	for _, it := range items {
		if it != nil && it.ID%100 == 0 {
			have[it.ID] = it.Title
		}
	}
	if have[actAutumnPrayerRoot] != "秋祈良愿" {
		t.Fatalf("autumn root title=%q", have[actAutumnPrayerRoot])
	}
	if have[actJoyShareRoot] != "快乐不独享" {
		t.Fatalf("joy root title=%q", have[actJoyShareRoot])
	}
}

func TestActEnsureGroupRootsKeepsExisting(t *testing.T) {
	items := actEnsureGroupRoots([]*ActivityInfo{
		{ID: actAutumnPrayerRoot, Title: "秋祈良愿"},
		{ID: actWishSignID, Title: "求签"},
	})
	n := 0
	for _, it := range items {
		if it != nil && it.ID == actAutumnPrayerRoot {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("root copies=%d", n)
	}
}

func TestActManualClaimAttemptsSkipHonghuaField(t *testing.T) {
	atts := actManualClaimAttempts(actJoyShareRoot+1, 21)
	if len(atts) == 0 || atts[0].extField != 21+honghuaExtBase+1 {
		t.Fatalf("first extField=%v want cmd+100", atts)
	}
	cmds := []int64{21, 4, 25, 1}
	for _, c := range cmds {
		if c == honghuaLoveCmd || c == honghuaFundCmd {
			t.Fatalf("probe still includes honghua cmd %d", c)
		}
	}
}
