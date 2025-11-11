package wizard_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWizardSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wizard Suite")
}
