package helpers

func FindBestMatch[S ~[]E, E any](x S, cmp func(a, b E) int) (*E, int) {
	if len(x) < 1 {
		panic("slices.MinFunc: empty list")
	}
	m := &x[0]
	ind := 0
	for i := 1; i < len(x); i++ {
		if cmp(x[i], *m) < 0 {
			m = &x[i]
			ind = i
		}
	}
	return m, ind
}
