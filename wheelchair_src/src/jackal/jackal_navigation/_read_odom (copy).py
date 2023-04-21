#!/usr/bin/env python3  
import rospy
import math
import tf
import geometry_msgs.msg
import time
from rosgraph_msgs.msg import Clock
from jackal_msgs.msg import Drive
from sensor_msgs.msg import JointState

sim_time = 0.000

def clock_callback(data):
    global sim_time
    sim_time = data.clock.secs + data.clock.nsecs*1000000
    
command = [0.0, 0.0]
def command_callback(data):
    global command
    command = [data.drivers[0], data.drivers[1]]
    
joints = [0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
def joints_callback(data):
    global joints
    joints = [data.position[0], data.position[1],data.position[2], data.position[3],
              data.velocity[0], data.velocity[1],data.velocity[2], data.velocity[3]]

def main():
    global sim_time, command, joints
    rospy.Subscriber("/clock", Clock, clock_callback)
    rospy.Subscriber("/cmd_drive", Drive, command_callback)
    rospy.Subscriber("/joint_states", JointState, joints_callback)
    listener = tf.TransformListener()
    rate = rospy.Rate(100.0)
    f = open('fast_record.csv', 'w')
    start_time = time.time()
    while not rospy.is_shutdown():
        current_time = time.time() - start_time
        try:
            try:
                (trans,rot) = listener.lookupTransform('/base_link', '/odom', rospy.Time(0))
            except (tf.LookupException, tf.ConnectivityException, tf.ExtrapolationException):
                continue
            rospy.loginfo(trans)
            data = trans
            data.extend(rot)
            data.extend(command)
            f.write(','.join('{:.5f}'.format(val) for val in data))
            f.write(',{:.5f}\n'.format(current_time))
            rate.sleep()
        except KeyboardInterrupt:
            print('recording fimished')
            f.close()
            break
    print('node stopped')
        
if __name__ == '__main__':
      rospy.init_node('recorder_node')
      try:
          main()
      except rospy.ROSInterruptException:
          pass
