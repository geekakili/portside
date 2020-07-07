package repository

// ArrayContains finds if an element (pin) is inside an array (hayStack)
func ArrayContains(hayStack []string, pin string) bool {
	for _, val := range hayStack {
		if val == pin {
			return true
		}
	}
	return false
}
