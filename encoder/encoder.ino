#include <EasyTransfer.h>
EasyTransfer ET;

const byte c1 = 2;
const byte c2 = 3;
const byte s1 = 19;
const byte s2 = 20;
const byte q1 = 17;
const byte q2 = 18;

bool isForward = true;
bool isForward2= true;

long tisks = 0;
long tisks2= 0;

struct SEND_DATA_STRUCTURE{
  unsigned long time;
  long tiks1;
  bool dir1;
  long tiks2;
  bool dir2;
  int  cast1;
  int  cast2;
};
SEND_DATA_STRUCTURE mydata;

void setup() {
  Serial.begin(115200);
  ET.begin(details(mydata), &Serial);

  pinMode(c1, INPUT);
  pinMode(c2, INPUT);
  pinMode(q1, INPUT);
  pinMode(q2, INPUT);
  pinMode(s1, INPUT);
  pinMode(s2, INPUT);
  
//  Serial.println("epta!");
  attachInterrupt(digitalPinToInterrupt(s1), trigg1, RISING);
  attachInterrupt(digitalPinToInterrupt(q2), trigg2, RISING);
}

void loop() {
  mydata.time = micros();
  mydata.tiks1= tisks;
  mydata.dir1 = isForward;
  mydata.tiks2= tisks2;
  mydata.dir2 = isForward;
  ET.sendData();
//  Serial.print(tisks);
//  Serial.print("\t");
//  Serial.println(tisks2);
}

void trigg1(){
//  Serial.print("TRIUGG");
  isForward= digitalRead(s2);
  if(isForward){
    tisks++;  
  }else{
    tisks--;  
  }
}

void trigg2(){
  isForward2= digitalRead(q1);
  if(isForward2){
    tisks2++;  
  }else{
    tisks2--;  
  }
}
