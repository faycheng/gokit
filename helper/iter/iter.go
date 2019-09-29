package iter

// Step .
type Step struct {
	Head int
	Tail int
}

// Range calculates the steps.
// for _, step := range pkg.RangeStep(len(mids), 10) {
// 		cur := mids[step.Head:step.Tail]
//	}
func Range(total, step int) (steps []Step) {
	return RangeStart(0, total, step)
}

// RangeStart calculates the steps.
func RangeStart(start, high, step int) (steps []Step) {
	steps = make([]Step, 0)
	for i := start; i < high; i++ {
		if i%step == 0 {
			head := i
			tail := head + step
			if tail > high {
				tail = high
			}
			steps = append(steps, Step{Head: head, Tail: tail})
		}
	}
	return steps
}
