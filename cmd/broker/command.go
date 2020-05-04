package broker

import (
	. "github.com/c1rno/idempotencer/cmd/shared"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/queue"
	_ "github.com/pebbe/zmq4"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   `broker`,
	Short: `0MQ broker, needs to load-balancing`,
	Run: func(cmd *cobra.Command, args []string) {
		setup, err := InitialSetup()
		helpers.Panicer(err)
		defer setup.Waiter()

		setup.Wg.Add(2)
		broker := queue.NewBroker(setup.Conf.QueueBroker, setup.Log)
		go func() {
			<-setup.Ctx.Done()
			helpers.Panicer(broker.Stop())
			setup.Wg.Done()
		}()
		go func() {
			err := broker.Start()
			if setup.Ctx.Err() == nil {
				helpers.Panicer(err)
			}
			setup.Wg.Done()
		}()
	},
}
