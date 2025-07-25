package iot_cli

import (
	"context"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func (c *SchedulerClient) SendMessage(topic string, payload []byte) error {
	var err error
	c.CliConfig, err = c.GetClientConfig()
	if err != nil {
		return err
	}
	ctx := context.Background()

	// Connect to the broker - this will return immediately after initiating the connection process
	cm, err := autopaho.NewConnection(ctx, *c.CliConfig)
	if err != nil {
		return err
	}

	err = cm.AwaitConnection(ctx)
	if err != nil {
		return err
	}
	go func(msg []byte) {
		pr, err := cm.Publish(ctx, &paho.Publish{
			QoS:     Qos,
			Topic:   topic,
			Payload: msg,
		})
		if err != nil {
			log.Debugf("error publishing: %s\n", err)
		} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
			log.Debugf("reason code %d received\n", pr.ReasonCode)
		} else {
			log.Debugf("sent message: %s\n", msg)
		}
	}(payload)

	select {
	case <-time.After(time.Second * 5):
	case <-ctx.Done():
		log.Debug("publisher done")
	}

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	disconnectCtx, disconnectCancel := context.WithTimeout(ctx, time.Second)
	defer disconnectCancel()
	_ = cm.Disconnect(disconnectCtx)

	log.Debug("shutdown complete")
	return nil
}
