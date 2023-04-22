v = 0; %1.11; % m/s which is 4km/h
% v = .83; % m/s which is 3km/h
w = 15;

L = 0.51; % m
R = 0.17; % m

v_goal = [(2*v - w*L) / (2*R) ; (2*v + w*L) / (2*R) ] % this is in rad/s

K1 = [ -18.5000  -18.5000;  5.4250   -5.4250];
K2 = [ -18.5000  -18.5000; 15.5000  -15.5000];
if (v_goal(1)+v_goal(2)>0)
    f_t = K2*v_goal;
else
    f_t = K1*v_goal;
end

f_t/2 + [171; 176]

