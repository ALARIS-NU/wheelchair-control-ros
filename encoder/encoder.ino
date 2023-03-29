//pin declaration
const byte c1 = 2;  /**/ //encoder/ 
const byte c2 = 3;  /**/ //-------/
const byte s1 = 9;

bool isForward = true;

long delta = 0;
long lastMessageMillis = 0;

long lastTrig = 0;
bool triggered = false;

long tisks = 0;

void setup() {
  delta = micros();
  lastMessageMillis = micros();
  Serial.begin(115200);
  pinMode(c1, INPUT);
  pinMode(c2, INPUT);
  pinMode(5, OUTPUT);
  pinMode(6, OUTPUT);
  pinMode(s1, INPUT);
  
  attachInterrupt(digitalPinToInterrupt(c1), trigg1, RISING);
  Serial.println("micros\ttikcs");
}

bool epta=true;
bool sync = true;

void loop() {
    if(sync){
      Serial.println(tisks);
    }else{
      sync=digitalRead(s1);
    }
}

void trigg1(){
  isForward= digitalRead(c2);
  if(isForward){
    tisks++;  
  }else{
    tisks--;  
  }
  triggered= true;
}
