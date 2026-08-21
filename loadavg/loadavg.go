// Package loadavg implements the exponentially smoothed load-average
// calculation used by the Linux kernel.
package loadavg

import (
	"errors"
	"math"
)

// Linux represents load averages as fixed-point values with 11 fractional
// bits. The EXP constants are the kernel's precomputed decay factors for a
// five-second update period.
const (
	FractionBits = 11
	FixedOne     = uint64(1 << FractionBits)
	Exp1         = uint64(1884)
	Exp5         = uint64(2014)
	Exp15        = uint64(2037)
)

var ErrActiveOverflow = errors.New("active task count is too large for fixed-point calculation")

// Averages contains Linux-compatible 1, 5, and 15 minute load averages.
// Its zero value is ready for use. Update must be called once per five-second
// sampling period, including periods where there are no active tasks.
type Averages struct {
	one, five, fifteen uint64
}

// Update incorporates the number of runnable or uninterruptible tasks in the
// current five-second sample. It performs the same fixed-point recurrence and
// upward rounding as Linux calc_load.
func (a *Averages) Update(active uint64) error {
	// The recurrence temporarily multiplies a Q11 value by another Q11
	// coefficient, so reserve enough room for that intermediate and rounding.
	maxActive := (uint64(math.MaxUint64) - (FixedOne - 1)) / (FixedOne * FixedOne)
	if active > maxActive {
		return ErrActiveOverflow
	}

	fixedActive := active * FixedOne
	a.one = calcLoad(a.one, Exp1, fixedActive)
	a.five = calcLoad(a.five, Exp5, fixedActive)
	a.fifteen = calcLoad(a.fifteen, Exp15, fixedActive)
	return nil
}

// Values returns the 1, 5, and 15 minute averages as floating-point numbers.
func (a Averages) Values() (one, five, fifteen float64) {
	return fromFixed(a.one), fromFixed(a.five), fromFixed(a.fifteen)
}

// FixedValues returns the raw Q11 values, which is useful when exact kernel
// compatibility matters.
func (a Averages) FixedValues() (one, five, fifteen uint64) {
	return a.one, a.five, a.fifteen
}

func calcLoad(load, exp, active uint64) uint64 {
	newLoad := load*exp + active*(FixedOne-exp)
	if active >= load {
		newLoad += FixedOne - 1
	}
	return newLoad / FixedOne
}

func fromFixed(value uint64) float64 {
	return float64(value) / float64(FixedOne)
}
