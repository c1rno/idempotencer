package downstream

import (
	"fmt"

	. "github.com/c1rno/idempotencer/cmd/shared"
	"github.com/c1rno/idempotencer/pkg/dto"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/queue"
	_ "github.com/pebbe/zmq4"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   `downstream`,
	Short: `Simple kafka consumer, produces events into 0MQ broker`,
	Run: func(cmd *cobra.Command, args []string) {
		setup, err := InitialSetup()
		helpers.Panicer(err)
		defer setup.Waiter()

		setup.Wg.Add(2)
		client := queue.NewClient(setup.Conf.QueueConsumer, setup.Log)
		go func() {
			<-setup.Ctx.Done()
			helpers.Panicer(client.Disconnect())
			setup.Wg.Done()
		}()
		go func() {
			var (
				err errors.Error
				msg dto.Msg
			)
			helpers.Panicer(client.Connect())
			err = client.Push(dto.NewRawMsg(queue.READY))
			for err != nil && setup.Ctx.Err() == nil {
				err = client.Push(dto.NewRawMsg(queue.READY))
			}
			id := helpers.UniqIdentity()
			i := 0
			for setup.Ctx.Err() == nil {
				i += 1
				msg, err = client.Pull()
				if err == nil {
					tmp := msg.Data()
					tmp[len(tmp)-1] = fmt.Sprintf("downstream-%s: %d", id, i)
					err = client.Push(dto.NewRawMsg(tmp...))
					for err != nil && setup.Ctx.Err() == nil {
						err = client.Push(dto.NewRawMsg(tmp...))
					}
				}
			}
			setup.Wg.Done()
		}()
	},
}
