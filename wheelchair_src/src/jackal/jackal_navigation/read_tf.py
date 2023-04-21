#!/usr/bin/env python3  
import rospy
import math
import tf
import geometry_msgs.msg
import time
from rosgraph_msgs.msg import Clock
from jackal_msgs.msg import Drive
from sensor_msgs.msg import JointState
from sensor_msgs.msg import Imu
from sensor_msgs.msg import MagneticField
from tf2_msgs.msg import TFMessage

sim_time = 0.000

def clock_callback(data):
    global sim_time
    sim_time = data.clock.secs + data.clock.nsecs*1000000
    
tf_array = []; # write time secs+nsecs*10e-9 + translation + rotation
def tf_callback(data):
    #print(len(data.transforms))
    global tf_array
    for element in data.transforms:
        if element.child_frame_id=='base_link':
            print(element.header.frame_id,element.transform.translation)
            tf_array.append([element.header.stamp.secs, element.header.stamp.nsecs, element.transform.translation.x, element.transform.translation.y, element.transform.translation.z, element.transform.rotation.x, element.transform.rotation.y, element.transform.rotation.z, element.transform.rotation.w])
    
command = [0.0, 0.0]
def command_callback(data):
    global command
    command = [data.drivers[0], data.drivers[1]]
    
joints = [0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
def joints_callback(data):
    global joints
    joints = [data.position[0], data.position[1],data.position[2], data.position[3],
              data.velocity[0], data.velocity[1],data.velocity[2], data.velocity[3]]


imu_data = [0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
def imu_callback(data):
    global imu_data
    imu_data = [data.orientation.x, data.orientation.y, data.orientation.z, data.orientation.w, data.angular_velocity.x,
              data.angular_velocity.y, data.angular_velocity.z, data.linear_acceleration.x, data.linear_acceleration.y, data.linear_acceleration.z]
              
zed_imu_data = [0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
def zed_imu_callback(data):
    global zed_imu_data
    zed_imu_data = [data.orientation.x, data.orientation.y, data.orientation.z, data.orientation.w, data.angular_velocity.x,
              data.angular_velocity.y, data.angular_velocity.z, data.linear_acceleration.x, data.linear_acceleration.y, data.linear_acceleration.z]

zed_mag_data = [0.0, 0.0, 0.0]
def zed_mag_callback(data):
    global zed_mag_data
    zed_mag_data = [data.magnetic_field.x, data.magnetic_field.y, data.magnetic_field.z]

def main():
    global sim_time, command, joints, imu_data, zed_imu_data,zed_mag_data, tf_array
    rospy.Subscriber("/clock", Clock, clock_callback)
    rospy.Subscriber("/cmd_drive", Drive, command_callback)
    rospy.Subscriber("/joint_states", JointState, joints_callback)
    rospy.Subscriber("/imu/data", Imu, imu_callback)
    rospy.Subscriber("/zed/zed_node/imu/data", Imu, zed_imu_callback)
    rospy.Subscriber("/zed/zed_node/imu/mag", MagneticField, zed_mag_callback)
    rospy.Subscriber("/tf", TFMessage, tf_callback)
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
            #rospy.loginfo(trans)
            data = trans
            data.extend(rot)
            data.extend(command)
            data.extend(joints)
            data.extend(imu_data)
            data.extend(zed_imu_data)
            data.extend(zed_mag_data)
            f.write(','.join('{:.5f}'.format(val) for val in data))
            f.write(',{:.5f}\n'.format(current_time))
            rate.sleep()
        except KeyboardInterrupt:
            print('recording fimished')
            f.close()
            break
    with open('tf_list.csv', 'w') as file:
        file.writelines('\t'.join(str(j) for j in i) + '\n' for i in tf_array)
    print('node stopped')
        
if __name__ == '__main__':
      rospy.init_node('recorder_node')
      try:
          main()
      except rospy.ROSInterruptException:
          pass
