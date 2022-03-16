package rob

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -write_package_comment=false -package=$GOPACKAGE -destination=mock_sim_test.go gitlab.com/akita/akita/v3/sim Port,Engine

func TestRob(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rob Suite")
}
