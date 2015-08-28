package pbc

func addSet(targets []string, str string) []string {
	found := false
	for _, target := range targets {
		if target == str {
			found = true
		}
	}

	if !found {
		return append(targets, str)
	}
	return targets
}
