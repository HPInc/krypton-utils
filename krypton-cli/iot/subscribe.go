package iot

import (
	"cli/cmd"
	"cli/iot/iot_cli"
	"flag"
	"fmt"
)

const (
	CmdSubscribeName = "subscribe"
	DeviceTasksTopic = "v1/%s/tasks"
)

// additional flags for subscribe
type SubscribeFlags struct {
	Reply *bool
}

type CmdSubscribe struct {
	cmd.CmdBase
	SchedulerBase
	LoginInput
	SubscribeFlags
}

func init() {
	commands[CmdSubscribeName] = NewCmdSubscribe()
}

func NewCmdSubscribe() *CmdSubscribe {
	c := CmdSubscribe{
		cmd.CmdBase{
			Name: CmdSubscribeName,
		},
		SchedulerBase{},
		LoginInput{},
		SubscribeFlags{},
	}
	fs := c.BaseInitFlags()
	(&c.SchedulerBase).initFlags(fs)
	(&c.SubscribeFlags).initFlags(fs)
	return &c
}

// initialize subscribe flags
func (f *SubscribeFlags) initFlags(fs *flag.FlagSet) {
	f.Reply = fs.Bool("reply", true, "reply to messages")
}

func (c *CmdSubscribe) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if c.Stdin {
		err = cmd.ErrParseStdin
	} else {
		c.DeviceTokenIn = *c.DeviceToken
		c.DeviceIdIn = *c.DeviceId
	}
	c.RunFunc = c.subscribe_device_topic
	return c, err
}

// subscribe to device topic
func (c *CmdSubscribe) subscribe_device_topic() {
	cli := iot_cli.NewClient(
		*c.IotServer,
		*c.CertFile,
		c.DeviceIdIn,
		c.DeviceTokenIn,
		*c.Timeout,
		*c.Protocol,
	)
	topic := fmt.Sprintf(DeviceTasksTopic, c.DeviceIdIn)
	// send message
	if err := cli.Subscribe(topic, *c.Reply); err != nil {
		log.Fatal("Login error: ", err)
	}
	fmt.Println("subscribe complete")
}

func (c *CmdSubscribe) GetInput() interface{} {
	return &c.LoginInput
}

func (c *CmdSubscribe) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}
