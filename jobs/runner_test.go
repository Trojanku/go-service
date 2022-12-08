package jobs_test

import (
	"Goo/integrationtest"
	"Goo/jobs"
	"Goo/model"
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

type testRegistry map[string]jobs.Func

func (r testRegistry) Register(name string, fn jobs.Func) {
	r[name] = fn
}

func TestRunner_Start(t *testing.T) {
	integrationtest.SkipIfShort(t)

	t.Run("starts the runner and runs jobs until the context is cancelled", func(t *testing.T) {
		queue, cleanup := integrationtest.CreateQueue()
		defer cleanup()

		log, logs := newLogger()

		runner := jobs.NewRunner(jobs.NewRunnerOptions{
			Log:   log,
			Queue: queue,
		})

		ctx, cancel := context.WithCancel(context.Background())

		runner.Register("test", func(ctx context.Context, m model.Message) error {
			foo, ok := m["foo"]
			require.True(t, ok)
			require.Equal(t, "bar", foo)

			cancel()
			return nil
		})

		err := queue.Send(context.Background(), model.Message{"job": "test", "foo": "bar"})
		require.NoError(t, err)

		// This blocks until the context is cancelled by the job function
		runner.Start(ctx)

		require.Equal(t, 3, logs.Len())
		require.Equal(t, "Starting", logs.All()[0].Message)
		require.Equal(t, "Successfully ran job", logs.All()[1].Message)
		require.Equal(t, "Stopping", logs.All()[2].Message)
	})
}

func newLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.InfoLevel)
	return zap.New(core), logs
}
