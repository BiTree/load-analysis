package loadavg

import (
	"errors"
	"math"
	"testing"
)

func TestKernelFixedPointSequence(t *testing.T) {
	var averages Averages
	want := [][3]uint64{
		{328, 68, 22},
		{630, 135, 44},
		{908, 201, 66},
	}

	for i, expected := range want {
		if err := averages.Update(2); err != nil {
			t.Fatalf("Update(2): %v", err)
		}
		one, five, fifteen := averages.FixedValues()
		if got := [3]uint64{one, five, fifteen}; got != expected {
			t.Fatalf("sample %d: got %v, want %v", i+1, got, expected)
		}
	}
}

func TestZeroValueAndDecay(t *testing.T) {
	var averages Averages
	one, five, fifteen := averages.Values()
	if one != 0 || five != 0 || fifteen != 0 {
		t.Fatalf("zero value returned %v, %v, %v", one, five, fifteen)
	}

	if err := averages.Update(4); err != nil {
		t.Fatal(err)
	}
	before, _, _ := averages.Values()
	if err := averages.Update(0); err != nil {
		t.Fatal(err)
	}
	after, _, _ := averages.Values()
	if after >= before {
		t.Fatalf("load did not decay: before=%v after=%v", before, after)
	}
}

func TestActiveOverflowDoesNotMutate(t *testing.T) {
	var averages Averages
	if err := averages.Update(1); err != nil {
		t.Fatal(err)
	}
	wantOne, wantFive, wantFifteen := averages.FixedValues()

	err := averages.Update(math.MaxUint64/FixedOne + 1)
	if !errors.Is(err, ErrActiveOverflow) {
		t.Fatalf("got error %v, want ErrActiveOverflow", err)
	}
	if one, five, fifteen := averages.FixedValues(); one != wantOne || five != wantFive || fifteen != wantFifteen {
		t.Fatal("overflowing update mutated averages")
	}
}
