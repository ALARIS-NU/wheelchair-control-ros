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
	v := msg.Linear.X
	w := msg.Angular.Z

	L := 0.51 // distance between wheels in meters
	R := 0.17 // wheel radius in meters
	WL := (2*v - w*L) / (2 * R)
	WR := (2*v + w*L) / (2 * R)
	fmt.Printf("WL: %+v, WR: %+v\n", WL, WR)

	v_goal := [2]float64{WL, WR}
	K1 := [2][2]float64{{-18.5000, -18.5000}, {5.4250, -5.4250}}
	K2 := [2][2]float64{{-18.5000, -18.5000}, {15.5000, -15.5000}}

	var f_t [2]float64

	if v_goal[0]+v_goal[1] > 0 {
		f_t[0] = K2[0][0]*v_goal[0] + K2[0][1]*v_goal[1]
		f_t[1] = K2[1][0]*v_goal[0] + K2[1][1]*v_goal[1]
	} else {
		f_t[0] = K1[0][0]*v_goal[0] + K1[0][1]*v_goal[1]
		f_t[1] = K1[1][0]*v_goal[0] + K1[1][1]*v_goal[1]
	}

	f_t[0] = f_t[0] + 171
	f_t[1] = f_t[1] + 176

	fmt.Println(f_t)

	var command commandPack
	command.Action = byte(CSetCh)
	EasyTransferSend(port, command)
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
