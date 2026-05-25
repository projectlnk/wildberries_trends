package window

import (
	"testing"
	"time"
)

func TestSlidingWindow_AddAndGet(t *testing.T) {
	win := NewSlidingWindow(2*time.Second, 1*time.Second)
	win.Add("a")
	win.Add("b")
	win.Add("a")
	counts := win.GetAllCounts()
	if counts["a"] != 2 || counts["b"] != 1 {
		t.Errorf("Expected a=2,b=1, got a=%d,b=%d", counts["a"], counts["b"])
	}
}

func TestSlidingWindow_Expiration(t *testing.T) {
	win := NewSlidingWindow(2*time.Second, 1*time.Second)
	win.Add("x")
	time.Sleep(2500 * time.Millisecond) // больше периода окна
	win.Add("y")
	counts := win.GetAllCounts()
	if _, ok := counts["x"]; ok {
		t.Error("x should have expired")
	}
	if counts["y"] != 1 {
		t.Errorf("Expected y=1, got %d", counts["y"])
	}
}
