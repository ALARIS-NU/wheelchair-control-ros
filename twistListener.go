package main

import (
	"fmt"
	"time"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

func on_twist(msg *geometry_msgs.Twist) {
	fmt.Printf("Incoming: %+v\n", msg)
	v := -msg.Linear.X
	w := msg.Angular.Z

	L := 0.51 // distance between wheels in meters
	R := 0.17 // wheel radius in meters
	WL := (2*v - w*L) / (2 * R)
	WR := (2*v + w*L) / (2 * R)

	// WL := v - w*L/2
	// WR := v + w*L/2

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

	f_t[0] = f_t[0]/2 + 171
	f_t[1] = -f_t[1]*8 + 176

	fmt.Println(f_t)

	send_joy(port, byte(f_t[0]), byte(f_t[1]))

}

func init_twistListener() {
	n, err := goroslib.NewNode(goroslib.NodeConf{
		// Name:          "goroslib_sub_twist",
		// Namespace:     "/wheelchair",
		// /wheelchair_velocity_controller/cmd_vel
		Name:          "cmd_vel",
		Namespace:     "/wheelchair_velocity_controller",
		MasterAddress: Arduino.rosMasterAdress,
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()
	// create a subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:  n,
		Topic: "cmd_vel",
		// Topic:    "wheelchair_twist_listener",
		Callback: on_twist,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	for {
		time.Sleep(1 * time.Hour)
	}
}
