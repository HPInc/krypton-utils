package ss

import (
	"bytes"
	"cli/cmd"
	"cli/common"
	"cli/config"
	"flag"
	"fmt"
	"io"
	"net/http"

	pb "cli/scheduler_protos"

	uuid "github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	CmdCreateTaskName = "create_task"
	ClaimTenantId     = "tid"
	ClaimDeviceId     = "sub"
	TaskRequestParam  = "task_request"
)

// additional flags for create_task
type CreateTaskFlags struct {
	Version       *uint
	ServiceId     *string
	Schedule      *string
	ConsignmentId *common.Uuid
	MessageType   *string
	Payload       *string
	TenantId      *string
	MessageId     *string
}

// take input from a device_token json as stdin
type LoginInput struct {
	DeviceTokenIn string `json:"device_token"`
	DeviceIdIn    string `json:"device_id"`
}

type CmdCreateTask struct {
	cmd.CmdBase
	SchedulerBase
	LoginInput
	CreateTaskFlags
}

func init() {
	commands[CmdCreateTaskName] = NewCmdCreateTask()
}

func NewCmdCreateTask() *CmdCreateTask {
	c := CmdCreateTask{
		cmd.CmdBase{
			Name: CmdCreateTaskName,
		},
		SchedulerBase{},
		LoginInput{},
		CreateTaskFlags{
			ConsignmentId: common.NewUUID(),
		},
	}
	fs := c.BaseInitFlags()
	(&c.SchedulerBase).initFlags(fs)
	(&c.CreateTaskFlags).initFlags(fs)
	return &c
}

func (f *CreateTaskFlags) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	f.Version = fs.Uint("version", 1, "message version")
	f.ServiceId = fs.String("service_id", s.GetManagementServer(), "task service id")
	f.TenantId = fs.String("tenant_id", s.GetTenantId(), "tenant id")
	f.Schedule = fs.String("schedule", "now", "task frequency")
	f.MessageType = fs.String("message_type", "CCS.GetConfig", "message type")
	f.Payload = fs.String("payload", `{"msg":"hello"}`, "payload")
	fs.Var(f.ConsignmentId, "consignment_id", "consignment id (uuid)")
	// uuid:uuid is the message id format expected by some clients so defaulting to it
	// this is a pass through field to end devices
	messageId := fmt.Sprintf("%v:%v", uuid.NewString(), uuid.NewString())
	f.MessageId = fs.String("message_id", messageId, "message id")
}

func (c *CmdCreateTask) verify() bool {
	if !c.ConsignmentId.IsSet() {
		c.ConsignmentId.SetDefault()
	}
	return *c.TenantId != "" &&
		*c.DeviceId != "" &&
		c.ConsignmentId.IsSet() &&
		*c.JwtToken != ""
}

func (c *CmdCreateTask) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if !c.verify() {
		if !c.Stdin {
			log.Error(
				"Please provide -device_id, -tenant_id, -consignment_id or specify -stdin for standard input")
			return nil, cmd.ErrMissingArgs
		} else {
			err = cmd.ErrParseStdin
		}
	}
	c.RunFunc = c.createTask
	return c, err
}

// get tenant id from cmdline param or from input token
func (c *CmdCreateTask) getTenantId() string {
	if *c.TenantId != "" {
		return *c.TenantId
	} else if c.DeviceTokenIn != "" {
		return getClaimFromToken(c.DeviceTokenIn, ClaimTenantId)
	}
	return ""
}

// get tenant id from cmdline param or from input token
func (c *CmdCreateTask) getDeviceId() string {
	if *c.DeviceId != "" {
		return *c.DeviceId
	} else if c.DeviceTokenIn != "" {
		return getClaimFromToken(c.DeviceTokenIn, ClaimDeviceId)
	}
	return ""
}

func (c *CmdCreateTask) createTask() {
	tenantId := c.getTenantId()
	deviceId := c.getDeviceId()

	taskRequest := &pb.CreateScheduledTaskRequest{
		Version:       uint32(*c.Version),
		ServiceId:     *c.ServiceId,
		DeviceIds:     []string{deviceId},
		ConsignmentId: c.ConsignmentId.String(),
		TenantId:      tenantId,
		Schedule:      *c.Schedule,
		MessageType:   *c.MessageType,
		Payload:       []byte(*c.Payload),
		MessageId:     *c.MessageId,
	}
	log.Debug("create_task request: ", taskRequest)

	requestBytes, err := proto.Marshal(taskRequest)
	if err != nil {
		log.Println("Request: ", taskRequest)
		log.Fatalf("Protobuf encode failed for task request! Error: %v\n", err)
	}
	log.Debug("create_task request (base64): ", common.GetBase64(requestBytes))

	req, err := http.NewRequest(http.MethodPost, c.getTaskUrl(),
		bytes.NewReader(requestBytes))
	if err != nil {
		log.Fatalf("failed to created new task request: %v\n", err)
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Authorization", "Bearer "+*c.JwtToken)
	req.Header.Add("Content-Type", "application/x-protobuf")

	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request: ", taskRequest)
		log.Fatalf("Task request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Println("Request: ", taskRequest)
		log.Fatalf("Task request failed with error code: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Request: ", taskRequest)
		log.Fatalf("Failed to read task response. Error: %v\n", err)
	}
	log.HttpResponse(resp, result)
	fmt.Println(string(result))
}

func (c *CmdCreateTask) GetInput() interface{} {
	return &c.LoginInput
}

func (c *CmdCreateTask) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

func (c *CmdCreateTask) getTaskUrl() string {
	return fmt.Sprintf("%s/%s/tasks", *c.HttpServer, *c.HttpBasePath)
}

// parse jwt token unverified and get claim if exists
func getClaimFromToken(token, claim string) string {
	val, err := common.GetClaimFromToken(token, claim)
	if err != nil {
		log.Fatalf("no claim %s in token", claim)
	}
	return val
}
