clear;
close;
load('sweep36.mat')

Data = table2array(sweep36);
time = Data(:,1)/1000000;
cmd_fwd = Data(:,2);
cmd_turn = Data(:,3);
vell_left = Data(:,4);
vell_right = Data(:,5);
%%
figure
plot(time, (cmd_fwd-174)/182*23)
hold on
plot(time, vell_left)
plot(time, -vell_right)
legend('cmd fow','left','right')
%%
figure
plot((cmd_fwd/255))
hold on
plot( vell_left)
plot( vell_right)
legend('cmd fow','left','right')
%%
figure
plot(time, cmd_fwd/255)
hold on
plot(time, cmd_turn/255)