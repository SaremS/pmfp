package scope

type ExactMatch struct {
	TargetHost string `json:"target_host"`
}

func (s *ExactMatch) IsInScope(target string) bool {
	return s.TargetHost == target
}

func init() {
	RegisterScope("exact_match", func() Scope { return &ExactMatch{} })
}
