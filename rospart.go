package main

import (
	"fmt"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
)

func initROS() {
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "wheelchair_remote_sub",
		MasterAddress: "192.168.0.220:11311",
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()

	// // create a subscriber
	// sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
	// 	Node:     n,
	// 	Topic:    "wheelchair_remote_topic",
	// 	Callback: onMessage,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer sub.Close()
}

func onMessage(msg *sensor_msgs.Imu) {
	fmt.Printf("Incoming: %+v\n", msg)
}
