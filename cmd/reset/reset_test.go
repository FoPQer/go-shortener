package main

import "testing"

type testResetter struct {
	value string
}

func (t *testResetter) Reset() {
	t.value = ""
}

func TestNewPoolInitializesMap(t *testing.T) {
	pool := NewPool[*testResetter, string]()
	if pool == nil {
		t.Fatal("expected pool to be initialized")
	}
	if pool.items == nil {
		t.Fatal("expected internal map to be initialized")
	}
}

func TestPoolGetReturnsZeroWhenMissing(t *testing.T) {
	pool := NewPool[*testResetter, string]()

	got := pool.Get("missing")
	if got != nil {
		t.Fatalf("expected zero value nil for missing key, got %#v", got)
	}
}

func TestPoolPutAndGetRoundTrip(t *testing.T) {
	pool := NewPool[*testResetter, int]()
	item := &testResetter{value: "payload"}

	pool.Put(42, item)
	got := pool.Get(42)

	if got != item {
		t.Fatalf("expected same pointer instance, got %#v", got)
	}
	if got.value != "payload" {
		t.Fatalf("expected value to be preserved, got %q", got.value)
	}
}
