package panictest

// PanicFunc demonstrates use of built-in panic — should be flagged.
func PanicFunc() {
	panic("something went wrong") // want `use of built-in panic is forbidden`
}

// PanicWithValue also uses built-in panic.
func PanicWithValue(err error) {
	if err != nil {
		panic(err) // want `use of built-in panic is forbidden`
	}
}
