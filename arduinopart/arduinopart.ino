/**
 * This is a Wheelchair control app made by ALARIS team in 2022
 * 
 * The voltage PWM values drop due to high current demand, 
 *  which alters the output voltage a bit accoarding to the table:
 * 0V   - 0V
 * 5V   - 3.654V
 * 2.5V - 1.830V
 * 
 * But the offset is handled by the frontend, so no worries.
 */

#include <EasyTransfer.h>
#define buttonPressTime 250
const bool DebugMode = true;
EasyTransfer ET; 

const byte ch1pins[] = {
   2, //back-front
   3, //left-right
};

const byte numOfChannels= sizeof ch1pins / sizeof ch1pins[0];
const byte pinOn   = 9;//on/off
const byte pinHorn  = 12;//horn
const byte pinUp   = 11;//speedUp
const byte pinDown = 10;//speedDown

enum    command{CUndefined, CTurnOn, CTurnOff,   CHorn, CSmaller, CBigger, CSetCh};
const byte pinMap[] = {255,   pinOn,    pinOn, pinHorn,  pinDown,   pinUp,    255};

struct RECEIVE_DATA_STRUCTURE{
  byte action;
  byte ch_name;
  byte value;
};

RECEIVE_DATA_STRUCTURE mydata;
bool isSending = false;
unsigned long lastSignalStart = millis();
byte theSendingPin = 255;

void setup(){
  Serial.begin(115200);
  Serial.println("hi!!!");
  ET.begin(details(mydata), &Serial);
  for (byte i=0; i<numOfChannels; i++){
    pinMode(    ch1pins[i], OUTPUT);
    analogWrite(ch1pins[i],    174); //174 cuz of voltage drop
  }
  for (byte i=CTurnOn; i<=CBigger; i++){
    pinMode(     pinMap[i], INPUT);
    digitalWrite(pinMap[i],   LOW);
  }
}

void triggerButton(byte pin){
  if (!isSending) {
    pinMode(pin,OUTPUT);
    digitalWrite(pin,LOW);
    theSendingPin = pin;
    lastSignalStart = millis(); 
    if (DebugMode) {Serial.print("Started pin "); Serial.println(pin);}
    isSending = true;
  }
}

void resetPinToDefault(){
  pinMode(theSendingPin,INPUT);
  if (DebugMode) { Serial.print("Ended pin "); Serial.println(theSendingPin); }
  theSendingPin = 255;
  isSending = false;
}

void setChannelOutput(byte pin, byte value){
  if (DebugMode) { 
    Serial.print("PWM on pin "); Serial.print(pin); 
    Serial.print(", with value "); Serial.println(value);
  }
  if (mydata.ch_name < numOfChannels)
    analogWrite(pin, value); //map(mydata.value, 0, 186, 0,255)
}

void loop(){
  if (isSending) {
    if (millis() > lastSignalStart + buttonPressTime) {
      resetPinToDefault();
    }
  }
  
  if(ET.receiveData()){
    switch (mydata.action){
      case CUndefined: break;
      case CTurnOn ... CBigger: {triggerButton(pinMap[mydata.action]); }break;
      case CSetCh: {setChannelOutput(ch1pins[mydata.ch_name], mydata.value); }break;
      default: break;
    }
  }
}
