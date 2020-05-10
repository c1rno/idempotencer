package nanomsg

import (
	"fmt"

	. "github.com/c1rno/idempotencer/cmd/shared"
	"github.com/c1rno/idempotencer/pkg/dto"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/queue"
	"github.com/spf13/cobra"
)

var DownstreamCmd = &cobra.Command{
	Use:   `downstream`,
	Short: `Consume events from Nanomsg broker`,
	Run: func(cmd *cobra.Command, args []string) {
		setup, err := InitialSetup()
		helpers.Panicer(err)
		defer setup.Waiter()

		setup.Wg.Add(2)
		client := queue.NewMangosClient(setup.Conf.QueueConsumer, setup.Log)
		go func() {
			<-setup.Ctx.Done()
			helpers.Panicer(client.Disconnect())
			setup.Wg.Done()
		}()
		go func() {
			var err errors.Error
			helpers.Panicer(client.Connect())
			err = client.Push(dto.NewStringMsg(queue.READY))
			for err != nil && setup.Ctx.Err() == nil {
				err = client.Push(dto.NewStringMsg(queue.READY))
			}
			id := helpers.UniqIdentity()
			i := 0
			for setup.Ctx.Err() == nil {
				i += 1
				_, err = client.Pull()
				if err == nil {
					msg := dto.NewByteMsg([]byte(fmt.Sprintf("downstream-%s: %d", id, i)))
					err = client.Push(msg)
					for err != nil && setup.Ctx.Err() == nil {
						err = client.Push(msg)
					}
				}
			}
			setup.Wg.Done()
		}()
	},
}
