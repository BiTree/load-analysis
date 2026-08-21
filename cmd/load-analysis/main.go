package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"load-analysis/loadavg"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "load-analysis:", err)
		os.Exit(1)
	}
}

func run() error {
	averages := loadavg.Averages{}
	in := bufio.NewReader(os.Stdin)
	for sample := 1; ; sample++ {
		var active uint64
		if _, err := fmt.Fscan(in, &active); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("read sample %d: %w", sample, err)
		}
		if err := averages.Update(active); err != nil {
			return fmt.Errorf("sample %d: %w", sample, err)
		}
		one, five, fifteen := averages.Values()
		fmt.Printf("%d\t%.2f %.2f %.2f\n", sample, one, five, fifteen)
	}
}
