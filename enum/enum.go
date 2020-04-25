package enum

// Policy in state
type Policy int8

// Policy const value
const (
	PolicyInvalid Policy = 0
	PolicyWrite   Policy = 1
	PolicyRead    Policy = 2
)
