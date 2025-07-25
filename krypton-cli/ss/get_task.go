package ss

import (
	"cli/cmd"
	"cli/common"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	CmdGetTaskName = "get_task"
	paramDeviceId  = "device_id"
)

type taskScheduled struct {
	DeviceId string `json:"device_id"`
	TaskId   string `json:"task_id"`
}

// mapping input from a create_task
//
//	{
//	 "version":1,
//	 "task_count":1,
//	 "consignment_id":"b02477b9-980d-453d-bd64-3836fd59646b",
//	 "tenant_id":"15d89f59-4e31-4bab-bea9-ecc6773eb67f",
//	 "tasks_scheduled":[
//	   {"task_id":"292c6cae-fd8b-11ed-9ac0-e26036419046","device_id":"8f3f9d93-3d23-4aa2-b648-a962277394ff","status":"queued"}
//	 ]
//	 }
type GetTaskInput struct {
	TasksScheduled []taskScheduled `json:"tasks_scheduled"`
}

type GetTaskFlags struct {
	TaskId *common.Uuid
}

type CmdGetTask struct {
	cmd.CmdBase
	SchedulerBase
	taskScheduled
}

var (
	getTaskFlags = GetTaskFlags{
		TaskId: common.NewUUID(),
	}
)

func init() {
	commands[CmdGetTaskName] = NewCmdGetTask()
}

func NewCmdGetTask() *CmdGetTask {
	c := CmdGetTask{
		cmd.CmdBase{
			Name: CmdGetTaskName,
		},
		SchedulerBase{},
		taskScheduled{},
	}
	fs := c.BaseInitFlags()
	(&c.SchedulerBase).initFlags(fs)
	(&getTaskFlags).initFlags(fs)
	return &c
}

// setup flags for interactive input via cmdline
func (f *GetTaskFlags) initFlags(fs *flag.FlagSet) {
	fs.Var(f.TaskId, "task_id", "task id (uuid)")
}

func (c *CmdGetTask) verify() bool {
	return getTaskFlags.TaskId.IsSet() && *c.JwtToken != ""
}

// input parse allows params and stdin
func (c *CmdGetTask) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if !c.verify() {
		if !c.Stdin {
			log.Error(
				"Please provide -task_id or specify -stdin for standard input")
			return nil, cmd.ErrMissingArgs
		} else {
			err = cmd.ErrParseStdin
		}
	} else {
		c.DeviceId = *c.SchedulerFlags.DeviceId
		c.TaskId = getTaskFlags.TaskId.String()
	}
	c.RunFunc = c.getTask
	return c, err
}

func (c *CmdGetTask) getTask() {
	data := url.Values{}
	data.Set(paramDeviceId, c.DeviceId)

	req, err := http.NewRequest(http.MethodGet, c.getTaskUrl(), nil)
	if err != nil {
		log.Fatalf("creating get task request failed: %v\n", err)
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Authorization", "Bearer "+*c.JwtToken)
	req.URL.RawQuery = data.Encode()

	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("get task failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("get task failed with error code: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read task response. Error: %v\n", err)
	}
	log.HttpResponse(resp, result)
	fmt.Println(string(result))
}

func (c *CmdGetTask) GetInput() interface{} {
	return &GetTaskInput{}
}

func (c *CmdGetTask) ExecuteWithArgs(i interface{}) error {
	var err error
	ti := i.(*GetTaskInput)
	for _, v := range ti.TasksScheduled {
		c.DeviceId = v.DeviceId
		c.TaskId = v.TaskId
		if err = c.Execute(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (c *CmdGetTask) getTaskUrl() string {
	return fmt.Sprintf("%s/%s/tasks/%s", *c.HttpServer, *c.HttpBasePath, c.TaskId)
}
