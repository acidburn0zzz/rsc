package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRsc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rsc Suite")
}
