package main

import (
	"fmt"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
)

func initROS() {
	node1, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "wheelchair_remote_sub",
		MasterAddress: "192.168.0.220:11311",
	})
	if err != nil {
		panic(err)
	}
	defer node1.Close()

	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     node1,
		Topic:    "wheelchair_move_command",
		Callback: onMessage,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()
}

func onMessage(msg *sensor_msgs.Imu) {
	fmt.Printf("Incoming: %+v\n", msg)
}
