package process

import (
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ContainerSuite{})
		s.AddSuite(&HealthSuite{})
		s.AddSuite(&ParallelInitializerSuite{})
		s.AddSuite(&RunnerSuite{})
		s.AddSuite(&WatcherSuite{})
	})
}
