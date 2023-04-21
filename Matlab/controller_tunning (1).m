load('sweep36.mat')
%%
% Comments by hand
% 175 - middle of the turn
% but no turn 176 and 0 speed 171 (precision error)
%%
Data = table2array(sweep36);
Data(432575:end,1)=Data(432575:end,1)+Data((432575-1),1);
time = Data(:,1)/1000000;
cmd_fwd = Data(:,2);
cmd_turn = Data(:,3);
vell_left = Data(:,4);
vell_right = Data(:,5);
%%
figure
hold on
plot(time, vell_left)
plot(time, vell_right)
plot(time, (cmd_turn-171)/(-37))
legend('left','right','cmd fow')
%%
% figure
% plot((cmd_fwd/255))
% hold on
% plot( vell_left)
% plot( vell_right)
% legend('cmd fow','left','right')
% %%
% figure
% plot(time, cmd_fwd/255)
% hold on
% plot(time, cmd_turn/255)
%%
% select unique commands
Commands = unique([cmd_fwd, cmd_turn],'rows');
new_data = [];
% get mean responce
for i = 1:length(Commands)
    x = Data((Data(:,2)==Commands(i,1)),:);
    x = x(x(:,3)==Commands(i,2),:);
    new_data(end+1,:) = [Commands(i,1), Commands(i,2), mean(x(:,4)), mean(x(:,5))];
end
%%
plot(new_data)
legend('fow','turn','left','right')
%%
figure
hold on
plot(new_data(:,3:4))
%plot((new_data(:,2)-171)/(-37))
legend('left','right')
%%
% plot((new_data(:,2)-171)/(-37)+(new_data(:,1)-176)/(-31))
% %%
% plot((new_data(:,2)-171)/(-37)+(new_data(:,1)-176)/(+31))
plot((new_data(:,2)-171)/(-37)+(new_data(:,1)-176)/(-0.35*31))
%%
% converson
output = [];
for i=1:length(new_data)
    if (new_data(i,1)>176)
        left = (new_data(i,2)-171)/(-37)+(new_data(i,1)-176)/(+0.35*31);
        right = (new_data(i,2)-171)/(-37)+(new_data(i,1)-176)/(-0.35*31);
        output(end+1,:) = [left, right];
    else
        left = (new_data(i,2)-171)/(-37)+(new_data(i,1)-176)/(+31);
        right = (new_data(i,2)-171)/(-37)+(new_data(i,1)-176)/(-31);
        output(end+1,:) = [left, right];
    end
end
%%
f_t = new_data(i,1:2)'-[171;176]
v = [1/(-37) 1/(+31); 1/(-37) 1/(-31);]*f_t
%%
% goal velocity
v_goal = [2.5; 1.5];
K1 = [ -18.5000  -18.5000;  5.4250   -5.4250];
K2 = [ -18.5000  -18.5000; 15.5000  -15.5000];
if (v_goal(1)+v_goal(2)>0)
    f_t = K2*v_goal;
else
    f_t = K1*v_goal;
end
% print command (forward, turn)
disp(f_t)
%%
plot(output(:,1:2))
legend('left','right','left model','right model')