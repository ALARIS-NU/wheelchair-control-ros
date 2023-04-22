# How to start ROS?
Start this in ALARIS-wheelchairControl/wheelchair_src
```Bash
roslaunch wheelchair_description spawn_wheelchair.launch 
```
Then start this is ALARIS-wheelchairControl/
```Bash
go run . --port /dev/ttyUSB0 --port_encoder /dev/ttyACM0 --gui=False
```
This is enough for now.

## pinout
```
const byte ch1pins[] = {
   2, //back-front
   3, //left-right
};

const byte pinOn   = 9;//on/off
const byte pinHorn  = 12;//horn
const byte pinUp   = 11;//speedUp
const byte pinDown = 10;//speedDown
```

## voltage
 * 0V   - 0V
 * 5V   - 3.654V
 * 2.5V - 1.830V