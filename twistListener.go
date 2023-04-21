package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

func on_twist(msg *geometry_msgs.Twist) {
	fmt.Printf("Incoming: %+v\n", msg)
}

func init_twistListener() {
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "goroslib_sub_twist",
		MasterAddress: Arduino.rosMasterAdress,
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()
	// create a subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     n,
		Topic:    "wheelchair_twist_listener",
		Callback: on_twist,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// wait for CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
