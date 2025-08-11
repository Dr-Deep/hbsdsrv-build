package triggers

import (
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
	"time"
)

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
		time.Sleep(time.Second * 5)

		c <- internal.TriggerSignal{
			JobName: t.Job,
			Reason:  "testing purpose",
		}
	}
}
