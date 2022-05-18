#include <EasyTransfer.h>
EasyTransfer ET; 
enum command{CUndefined,CTurnOn, CTurnOff, CHorn, CSmaller, CBigger, CSetCh};

const byte ch1pins[] = {
   13,//channel1
   12,//channel2
   11,//channel3
   10 //channel4
};

struct RECEIVE_DATA_STRUCTURE{
  byte action;
  byte ch_name;
  byte value;
};

RECEIVE_DATA_STRUCTURE mydata;

void setup(){
  Serial.begin(115200);
  ET.begin(details(mydata), &Serial);
  pinMode(13, OUTPUT);
}

void loop(){
  if(ET.receiveData()){
    switch (mydata.action){
      case CUndefined: break;
      case CTurnOff: digitalWrite(13,LOW); break;
      case CTurnOn:  digitalWrite(13,HIGH); break;
      case CSetCh: {
        analogWrite(ch1pins[mydata.ch_name], mydata.value);
      }break;
      default: break;  
    }
  }
}
