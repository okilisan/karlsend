package txscript

import (
	"os"
	"testing"

	"github.com/karlsen-network/karlsend/v2/infrastructure/logger"
)

func TestMain(m *testing.M) {
	// set log level to trace, so that logClosures passed to log.Tracef are covered
	log.SetLevel(logger.LevelTrace)
	logger.InitLogStdout(logger.LevelTrace)

	os.Exit(m.Run())
}
