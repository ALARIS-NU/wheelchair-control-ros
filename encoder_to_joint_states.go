package main

import (
	"fmt"
	"math"
	"time"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type EncoderToJointStates struct {
	node                      *goroslib.Node
	jointStatesPub            *goroslib.Publisher
	wheelRadius               float64
	encoderTicksPerRevolution int
}

func NewEncoderToJointStates() (*EncoderToJointStates, error) {
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:      "encoder_to_joint_states",
		Namespace: "/wheelchair",
	})
	if err != nil {
		return nil, err
	}

	jsPub, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  n,
		Topic: "joint_states",
		Msg:   &sensor_msgs.JointState{},
	})
	if err != nil {
		return nil, err
	}

	e := &EncoderToJointStates{
		node:                      n,
		jointStatesPub:            jsPub,
		wheelRadius:               0.17, // 17cm radius
		encoderTicksPerRevolution: 2048,
	}

	return e, nil
}

func (e *EncoderToJointStates) publishJointStates() {
	jointStates := &sensor_msgs.JointState{
		Header: std_msgs.Header{
			Stamp: time.Now(),
		},
		Name: []string{"left_wheel_joint", "right_wheel_joint"},
	}

	leftWheelAngle := (float64(left_ticks) / float64(e.encoderTicksPerRevolution)) * (2 * math.Pi)
	rightWheelAngle := (float64(right_ticks) / float64(e.encoderTicksPerRevolution)) * (2 * math.Pi)

	jointStates.Position = []float64{leftWheelAngle, rightWheelAngle}

	e.jointStatesPub.Write(jointStates)
}

func init_encoder_joint_states() {
	e, err := NewEncoderToJointStates()
	if err != nil {
		fmt.Printf("error while initializing encoder to joint states node: %v", err)
		return
	}

	defer e.node.Close()

	for {
		e.publishJointStates()
		time.Sleep(10 * time.Millisecond) // Adjust the sleep duration according to your desired update rate
	}
}
