package main

import (
	"regexp"
	"net/http"
)

type ScoperChecker interface {
	IsInScope(target string) bool
}

type ExactMatchScope {
	TargetHost string `yaml:"target_host"`
} 

func (s *ExactMatchScope) IsInScope(target string) bool {
	return s.TargetHost == target
}

type RequestManipulator interface {
	ManipulateRequest(request http.Request) http.Request
}

type Config struct {
	Scope []ScopeChecker `yaml:"scope_checker"`
	Manipulators []RequestManipulators `yaml:"request_manipulators"`
}

