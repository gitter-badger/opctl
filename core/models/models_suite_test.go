package models

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestModels(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "engine/core/models")
}
