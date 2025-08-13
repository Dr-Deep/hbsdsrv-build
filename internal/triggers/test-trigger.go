package triggers

import (
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
	"time"
)

const testIntervalDuration = time.Second * 5

type TriggerTest struct {
	config.TriggerConfig
}

func NewTriggerTest(cfg config.TriggerConfig) internal.Trigger {
	return &TriggerTest{TriggerConfig: cfg}
}

func (t *TriggerTest) Run(b *internal.Builder, c chan internal.TriggerSignal) {

	b.Logger.Debug("running")
	// solang channel auf ist

	for {
		time.Sleep(testIntervalDuration)

		c <- internal.TriggerSignal{
			JobName: t.Job,
			Reason:  "testing purpose",
		}
	}
}
