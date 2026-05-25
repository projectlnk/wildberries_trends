package stoplist

import "testing"

func TestStopList_AddContains(t *testing.T) {
	sl := New()
	sl.Add("badword")
	if !sl.Contains("badword") {
		t.Error("Contains should return true for added word")
	}
	if sl.Contains("good") {
		t.Error("Contains should return false for missing word")
	}
}

func TestStopList_Remove(t *testing.T) {
	sl := New()
	sl.Add("bad")
	sl.Remove("bad")
	if sl.Contains("bad") {
		t.Error("Word should be removed")
	}
}

func TestStopList_List(t *testing.T) {
	sl := New()
	sl.Add("one")
	sl.Add("two")
	list := sl.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 words, got %d", len(list))
	}
}
