package trigger

// Trigger represents an abstract draining trigger
type Trigger interface {
	Loop()
}
