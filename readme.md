# How to start ROS?
Start this in ALARIS-wheelchairControl/wheelchair_src
```Bash
roslaunch wheelchair_description spawn_wheelchair.launch 
```
Then start this is ALARIS-wheelchairControl/
```Bash
go run . --port /dev/ttyUSB0 --port_encoder /dev/ttyACM0 --gui=False
```

Start the ZED camera
```Bash
roslaunch zed_wrapper zed.launch
```

The movebase is currently resides in
```bash
roslaunch zed_rtabmap_example jackal_robot_rtabmap.launch lidar2d:=false lidar3d:=false
```


This should give you two rviz tabs, first one is redundant but can be used with old jackal navigation to use only odom.

New rviz takes some time to start the camera, but then runs good. Change its 2D map topic to costmap to use NavGoal

# TODO
Current goal is to shift costmap plane in camera and change camera angle

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