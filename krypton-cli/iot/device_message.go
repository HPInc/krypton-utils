package iot

import (
	"cli/cmd"
	"cli/common"
	"cli/iot/iot_cli"
	"flag"
	"fmt"

	pb "cli/scheduler_protos"

	"google.golang.org/protobuf/proto"
)

const (
	CmdDeviceMessageName = "device_message"
)

// additional flags for device_message
type DeviceMessageFlags struct {
	Version     *uint
	MessageType *string
	Payload     *string
	Topic       *string
}

type CmdDeviceMessage struct {
	cmd.CmdBase
	SchedulerBase
	LoginInput
	DeviceMessageFlags
}

type TaskDetails struct {
	Success bool `json:"success"`
}

func init() {
	commands[CmdDeviceMessageName] = NewCmdDeviceMessage()
}

func NewCmdDeviceMessage() *CmdDeviceMessage {
	c := CmdDeviceMessage{
		cmd.CmdBase{
			Name: CmdDeviceMessageName,
		},
		SchedulerBase{},
		LoginInput{},
		DeviceMessageFlags{},
	}
	fs := c.BaseInitFlags()
	(&c.SchedulerBase).initFlags(fs)
	(&c.DeviceMessageFlags).initFlags(fs)
	return &c
}

// initialize device_message flags
func (f *DeviceMessageFlags) initFlags(fs *flag.FlagSet) {
	f.Version = fs.Uint("version", 1, "message version")
	f.MessageType = fs.String("message_type", "msg", "message type")
	f.Payload = fs.String("payload", `{"msg":"hello"}`, "payload")
	f.Topic = fs.String("topic", `v1/@cloud`, "send message to this topic")
}

func (c *CmdDeviceMessage) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if c.Stdin {
		err = cmd.ErrParseStdin
	} else {
		c.DeviceTokenIn = *c.DeviceToken
		c.DeviceIdIn = *c.DeviceId
	}
	c.RunFunc = c.send_device_message
	return c, err
}

// send message to server
func (c *CmdDeviceMessage) send_device_message() {
	cli := iot_cli.NewClient(
		*c.IotServer,
		*c.CertFile,
		c.DeviceIdIn,
		c.DeviceTokenIn,
		*c.Timeout,
		*c.Protocol,
	)
	message := c.getProtoMessage()
	// send message
	if err := cli.SendMessage(*c.Topic, message); err != nil {
		log.Fatal("Login error: ", err)
	}
	fmt.Println(common.GetJsonString(TaskDetails{Success: true}))
}

// encode a protobuf message to send to server
func (c *CmdDeviceMessage) getProtoMessage() []byte {
	msg, err := proto.Marshal(&pb.DeviceMessage{
		Version:     1,
		AccessToken: c.DeviceTokenIn,
		MessageType: *c.MessageType,
		Payload:     []byte(*c.Payload),
	})
	if err != nil {
		log.Fatal("Protobuf message encode error: ", err)
	}
	return msg
}

func (c *CmdDeviceMessage) GetInput() interface{} {
	return &c.LoginInput
}

func (c *CmdDeviceMessage) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}
