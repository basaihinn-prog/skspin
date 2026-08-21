package linedirection

import "fmt"

type T uint8

const (
	Left2Right T = iota
	Right2Left
	Both
)

func (this T) String() string {
	switch this {
	case Left2Right:
		return "从左到右"
	case Right2Left:
		return "从右到左"
	case Both:
		return "双向"
	default:
		return fmt.Sprintf("未知(%d)", this)
	}
}
