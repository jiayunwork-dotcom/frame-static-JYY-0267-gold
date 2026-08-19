package model

// Validate checks structural soundness and returns an *InvalidModelError for
// the first defect that would make the stiffness problem ill-posed. For the
// full list of defects, use ValidateAll.
func Validate(m *Model) error {
	errs := ValidateAll(m)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
