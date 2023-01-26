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
   5, //back-front
   3, //left-right
};

const byte numOfChannels= sizeof ch1pins / sizeof ch1pins[0];
const byte pinOn   = 9;//on/off
const byte pinHorn = 12;//horn
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
long reference = 5000;
float correction = 1.0f;

void setup(){
  Serial.begin(115200);
  
  //set the D6 to 2.5V with internal reference
  analogReference(EXTERNAL);
  reference = readVcc();
  Serial.println(reference);
  int initailV = 870000L/reference;
  analogWrite(6,640000L/reference);
  correction = 5000.0/float(reference);
  
  if (DebugMode) { Serial.println("hi!!!"); }
  ET.begin(details(mydata), &Serial);
  for (byte i=0; i<numOfChannels; i++){
    pinMode(    ch1pins[i], OUTPUT);
    Serial.print("setting port to..");
    Serial.println(initailV);
    analogWrite(ch1pins[i], initailV); //174 cuz of voltage drop
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
    analogWrite(pin, int(float(value)*correction)); //map(mydata.value, 0, 186, 0,255)
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


long readVcc() {
  // Read 1.1V reference against AVcc
  // set the reference to Vcc and the measurement to the internal 1.1V reference
  #if defined(__AVR_ATmega32U4__) || defined(__AVR_ATmega1280__) || defined(__AVR_ATmega2560__)
    ADMUX = _BV(REFS0) | _BV(MUX4) | _BV(MUX3) | _BV(MUX2) | _BV(MUX1);
  #elif defined (__AVR_ATtiny24__) || defined(__AVR_ATtiny44__) || defined(__AVR_ATtiny84__)
    ADMUX = _BV(MUX5) | _BV(MUX0);
  #elif defined (__AVR_ATtiny25__) || defined(__AVR_ATtiny45__) || defined(__AVR_ATtiny85__)
    ADMUX = _BV(MUX3) | _BV(MUX2);
  #else
    ADMUX = _BV(REFS0) | _BV(MUX3) | _BV(MUX2) | _BV(MUX1);
  #endif  

  delay(2); // Wait for Vref to settle
  ADCSRA |= _BV(ADSC); // Start conversion
  while (bit_is_set(ADCSRA,ADSC)); // measuring

  uint8_t low  = ADCL; // must read ADCL first - it then locks ADCH  
  uint8_t high = ADCH; // unlocks both

  long result = (high<<8) | low;

  result = 1086364L / result; // Calculate Vcc (in mV); 1125300 = 1.1*1023*1000
  return result; // Vcc in millivolts
}
