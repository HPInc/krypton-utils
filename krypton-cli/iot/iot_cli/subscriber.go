package iot_cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "cli/scheduler_protos"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"

	"google.golang.org/protobuf/proto"
)

type subscriber struct {
	cli         *SchedulerClient
	shouldReply bool
	cm          *autopaho.ConnectionManager
}

func NewSubscriber(cli *SchedulerClient, shouldReply bool) *subscriber {
	return &subscriber{cli: cli, shouldReply: shouldReply}
}

func (s *subscriber) Close() {
}

func (s *subscriber) setConnection(cm *autopaho.ConnectionManager) {
	s.cm = cm
}

func (s *subscriber) handle(msg *paho.Publish) {
	log.Printf("received message: %s\n", msg.Payload)
	if s.shouldReply {
		if err := s.publishReply(msg.Payload); err != nil {
			log.Printf("Reply failed: %v\n", err)
		}
		log.Printf("replied to message")
	}
}

// publish using the same client
func (s *subscriber) publishReply(svcPayload []byte) error {
	var err error

	replyPayload, err := s.makeReplyPayload(svcPayload)
	if err != nil {
		return err
	}

	if s.cm != nil {
		_, err = s.cm.Publish(context.Background(), &paho.Publish{
			QoS:     Qos,
			Topic:   "v1/@cloud/task_responses",
			Payload: replyPayload,
		})
	}
	return err
}

// make a device payload from received service payload
func (s *subscriber) makeReplyPayload(svcPayload []byte) ([]byte, error) {
	var svcMsg pb.ServiceMessage
	if err := proto.Unmarshal(svcPayload, &svcMsg); err != nil {
		log.Println("svc payload unmarshal error: ", err)
		return nil, err
	}

	log.Println("Reply: ", svcMsg.MessageType, svcMsg.TaskId)

	msg, err := proto.Marshal(&pb.DeviceMessage{
		Version:     1,
		AccessToken: s.cli.Token,
		MessageType: svcMsg.MessageType,
		TaskId:      svcMsg.TaskId,
		TaskStatus:  "complete",
		Payload:     svcMsg.Payload,
	})
	if err != nil {
		log.Println("Protobuf message encode error: ", err)
		return nil, err
	}
	return msg, nil
}

// subscribe
// use the default config, then override subscription handlers
// - subscribe on connection up,
// - route messages to subscription handler
func (c *SchedulerClient) Subscribe(topic string, shouldReply bool) error {
	var err error
	if c.CliConfig, err = c.GetClientConfig(); err != nil {
		return err
	}

	ctx := context.Background()

	subscriber := NewSubscriber(c, shouldReply)
	defer subscriber.Close()

	c.Topic = topic
	c.CliConfig.OnConnectionUp = c.SubscribeOnConnectionUp
	c.CliConfig.ClientConfig.Router = paho.NewSingleHandlerRouter(func(m *paho.Publish) {
		subscriber.handle(m)
	})

	// Connect to the broker - this will return immediately after initiating the connection process
	cm, err := autopaho.NewConnection(ctx, *c.CliConfig)
	if err != nil {
		return err
	}

	subscriber.setConnection(cm)

	err = cm.AwaitConnection(ctx)
	if err != nil {
		return err
	}

	// Messages will be handled through the callback so we really just need to wait until a shutdown
	// is requested
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	log.Println("signal caught - exiting")

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	disconnectCtx, disconnectCancel := context.WithTimeout(ctx, time.Second)
	defer disconnectCancel()
	_ = cm.Disconnect(disconnectCtx)

	log.Debug("shutdown complete")
	return nil
}

// start subscription on connection up
func (c *SchedulerClient) SubscribeOnConnectionUp(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
	if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			c.Topic: {QoS: Qos},
		},
	}); err != nil {
		log.Printf("failed to subscribe (%s). no messages will be received.", err)
		return
	}
	log.Println("mqtt subscribed to: ", c.Topic)
}
