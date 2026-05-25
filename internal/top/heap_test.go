package top

import "testing"

func TestGetTopN(t *testing.T) {
	counts := map[string]int64{
		"a": 10,
		"b": 5,
		"c": 20,
		"d": 1,
	}
	top3 := GetTopN(counts, 3)
	expected := []string{"c", "a", "b"}
	if len(top3) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(top3))
	}
	for i, v := range expected {
		if top3[i] != v {
			t.Errorf("Position %d: expected %s, got %s", i, v, top3[i])
		}
	}
	top10 := GetTopN(counts, 10)
	if len(top10) != 4 {
		t.Errorf("Expected 4 items, got %d", len(top10))
	}
}

func TestGetTopN_Empty(t *testing.T) {
	top := GetTopN(map[string]int64{}, 5)
	if len(top) != 0 {
		t.Error("Expected empty slice")
	}
}
