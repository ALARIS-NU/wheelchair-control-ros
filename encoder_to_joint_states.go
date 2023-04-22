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
	prevLeftTicks             int
	prevRightTicks            int
	prevTimestamp             time.Time
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
	now := time.Now()

	leftWheelAngle := (float64(left_ticks) / float64(e.encoderTicksPerRevolution)) * (2 * math.Pi)
	rightWheelAngle := (float64(right_ticks) / float64(e.encoderTicksPerRevolution)) * (2 * math.Pi)

	deltaTime := now.Sub(e.prevTimestamp).Seconds()
	leftWheelVelocity := ((float64(left_ticks) - float64(e.prevLeftTicks)) / float64(e.encoderTicksPerRevolution)) * (2 * math.Pi) / deltaTime
	rightWheelVelocity := ((float64(right_ticks) - float64(e.prevRightTicks)) / float64(e.encoderTicksPerRevolution)) * (2 * math.Pi) / deltaTime

	e.prevLeftTicks = left_ticks
	e.prevRightTicks = right_ticks
	e.prevTimestamp = now

	jointStates := &sensor_msgs.JointState{
		Header: std_msgs.Header{
			Stamp: now,
		},
		// Name:     []string{"Rev12", "Rev15"},
		Name:     []string{"Rev15", "Rev12"},
		Position: []float64{leftWheelAngle, rightWheelAngle},
		Velocity: []float64{leftWheelVelocity, rightWheelVelocity},
	}

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
