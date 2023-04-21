package main

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
	"github.com/fatih/color"
)

const logEncoder = true

const (
	preambleFirst  byte = 0x06
	preambleSecond byte = 0x85
)

type SendDataStructure struct {
	Time  uint32
	Tiks1 int32
	Dir1  bool
	Tiks2 int32
	Dir2  bool
	Cast1 int16
	Cast2 int16
}

var left_ticks int
var right_ticks int
var micros uint32
var prev_left_ticks int
var prev_right_ticks int
var prev_micros uint32
var left_speed float32
var right_speed float32

var writer *csv.Writer

func read_encoder(f io.ReadWriteCloser) {
	if Arduino.log_speed {
		t := time.Now()
		filename := fmt.Sprintf("SpeedLog-%d-%02d-%02d_%02d-%02d-%02d.csv", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

		file, err := os.Create(filename)
		if err != nil {
			log.Fatal("Error creating file:", err)
		}
		defer file.Close()

		// Create a new CSV writer
		writer = csv.NewWriter(file)
	}
	color.Green("read_encoder start")
	reader := bufio.NewReader(f)
	for {
		packet, err := readPacket(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading packet:", err)
			}
			// break
		} else {

			if packet.Time-micros > 5000 {
				prev_micros = micros
				prev_left_ticks = left_ticks
				prev_right_ticks = right_ticks

				left_ticks = int(packet.Tiks1)
				right_ticks = int(packet.Tiks2)
				micros = packet.Time

				// left and right ticks in radians per second if one revolution is 1024 ticks
				left_speed = float32(left_ticks-prev_left_ticks) * 2 * 3.14159 / 2048 / float32(micros-prev_micros) * 1000000
				right_speed = float32(right_ticks-prev_right_ticks) * 2 * 3.14159 / 2048 / float32(micros-prev_micros) * 1000000

				if Arduino.log_speed {
					data := [][]string{
						{"Time", "forward", "right", "left_speed", "right_speed"},
						{strconv.FormatUint(uint64(micros), 10), strconv.FormatUint(uint64(Arduino.forward), 10), strconv.FormatUint(uint64(Arduino.right), 10), strconv.FormatFloat(float64(left_speed), 'f', -1, 32), strconv.FormatFloat(float64(right_speed), 'f', -1, 32)},
					}
					for _, row := range data {
						err := writer.Write(row)
						if err != nil {
							log.Fatal("Error writing row:", err)
						}
					}

					writer.Flush()

					if err := writer.Error(); err != nil {
						log.Fatal("Error flushing writer:", err)
					}
				}
				// fmt.Println("Received packet:", packet)
				// fmt.Printf("left_speed: %2.2f, right_speed: %2.2f", left_speed, right_speed)
			}
		}
	}
}

func readPacket(reader *bufio.Reader) (SendDataStructure, error) {
	var data SendDataStructure

	var prev byte
	for {
		// Read until we find the two-byte preamble
		b, err := reader.ReadByte()
		if err != nil {
			return data, err
		}

		if prev == preambleFirst && b == preambleSecond {
			break
		}
		prev = b
	}

	// Read the length of the payload
	_, err := reader.ReadByte()
	// color.Green("length %v", length)
	if err != nil {
		return data, err
	}

	// Read the payload
	payload := make([]byte, 18)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return data, err
	}

	// Read the checksum
	checksum, err := reader.ReadByte()
	if err != nil {
		return data, err
	}

	// Validate the checksum
	calculatedChecksum := calculateChecksum(payload)
	if calculatedChecksum != checksum {
		return data, fmt.Errorf("invalid checksum: expected %02x, got %02x", calculatedChecksum, checksum)
		// color.Red("invalid checksum: expected %02x, got %02x", calculatedChecksum, checksum)

	}

	// Parse the payload bytes into the SendDataStructure struct
	data.Time = binary.LittleEndian.Uint32(payload[0:4])
	data.Tiks1 = int32(binary.LittleEndian.Uint32(payload[4:8]))
	data.Dir1 = payload[8] != 0
	data.Tiks2 = int32(binary.LittleEndian.Uint32(payload[9:13]))
	data.Dir2 = payload[13] != 0
	data.Cast1 = int16(binary.LittleEndian.Uint16(payload[14:16]))
	data.Cast2 = int16(binary.LittleEndian.Uint16(payload[16:18]))

	return data, nil
}

func calculateChecksum(data []byte) byte {
	// color.Blue("calculateChecksum %v", data)
	// data = append([]byte{0x34}, data...)
	var checksum byte
	checksum = 0x12
	for _, b := range data {
		checksum ^= b
	}
	return checksum
}

func init_encoder_ROS() {
	// create a node and connect to the master
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "goroslib_pub_encoder",
		Namespace:     "/wheelchair/",
		MasterAddress: Arduino.rosMasterAdress,
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()

	// create a publisher
	pub, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  n,
		Topic: "encoder_rad",
		Msg:   &std_msgs.Float32MultiArray{},
	})
	if err != nil {
		panic(err)
	}
	defer pub.Close()

	// Topic list
	// wheel_left
	// wheel_right
	// Rad/sec (single topic) FloatMultiArray[left, right];

	// r := n.TimeRate(1 * time.Second)

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)

	for {
		// select {
		// publish a message every second
		// case <-r.SleepChan():
		msg := &std_msgs.Float32MultiArray{
			Data: []float32{left_speed, -1.0 * right_speed},
		}
		// fmt.Printf("Outgoing: %+v\n", msg)
		pub.Write(msg)

		// handle CTRL-C
		// case <-c:
		// 	return
		// }
	}
}
