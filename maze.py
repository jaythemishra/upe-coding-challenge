#!/usr/bin/python
import requests, json, sys

#Written by Jay Mishra
UID = '704925466'
URL = "http://ec2-34-216-8-43.us-west-2.compute.amazonaws.com"
NEWGAMEURL = URL + "/session"
r = requests.post(NEWGAMEURL, data = {'uid':UID})
converted = json.loads(r.text)
TOKEN = converted['token']
GAMEURL = URL + "/game?token=" + TOKEN

r = requests.get(GAMEURL)
converted = json.loads(r.text)
status = converted['status']
print("Status: " + status)
if status == 'NONE' or status == 'FINISHED':
    sys.exit()
x = converted['current_location'][0]
y = converted['current_location'][1]
totalX = converted['maze_size'][0]
totalY = converted['maze_size'][1]
levels_completed = converted['levels_completed']
total_levels = converted['total_levels']

moves = []
maze = [[] for i in range(totalX)]

for i in range(totalX):
    for j in range (totalY):
        maze[i].append('_');

print('AND WE\'RE OFF!')
print()

while (1):
    backtracked = False
    r = requests.get(GAMEURL)
    gameStatus = json.loads(r.text)
    status = gameStatus['status']
    if (status == 'NONE'):
        print("Session has expired or does not exists")
        sys.exit()
    elif (status == 'FINISHED'):
        print("All levels completed")
        sys.exit()
    x = gameStatus['current_location'][0]
    y = gameStatus['current_location'][1]
    totalX = gameStatus['maze_size'][0]
    totalY = gameStatus['maze_size'][1]
    levels_completed = gameStatus['levels_completed']
    total_levels = gameStatus['total_levels']
    move = 'LEFT'
    if x > 0 and maze[x-1][y] != 'X':
        move = 'LEFT'
    elif x < totalX - 1 and maze[x+1][y] != 'X':
        move = 'RIGHT'
    elif y > 0 and maze[x][y-1] != 'X':
        move = 'UP'
    elif y < totalY - 1 and maze[x][y+1] != 'X':
        move = 'DOWN'
    else:
        move = moves.pop()
        maze[x][y] = 'X'
        print("Backtracking " + move)
        backtracked = True
    r = requests.post(GAMEURL, data = {'action':move})
    results = json.loads(r.text)
    result = results['result']
    if result == 'END':
        print()
        print("Level " + str(levels_completed + 1) + " of " + str(total_levels) + " completed.")
        print()
        r = requests.get(GAMEURL)
        converted = json.loads(r.text)
        status = converted['status']
        if (status == 'NONE'):
            print("Session has expired or does not exists")
            sys.exit()
        elif (status == 'FINISHED'):
            print("All levels completed")
            sys.exit()
        x = converted['current_location'][0]
        y = converted['current_location'][1]
        totalX = converted['maze_size'][0]
        totalY = converted['maze_size'][1]
        levels_completed = converted['levels_completed']
        total_levels = converted['total_levels']
        maze = [[] for i in range(totalX)]
        for i in range(totalX):
            for j in range (totalY):
                maze[i].append('_');
    elif result == 'WALL' or result == 'OUT_OF_BOUNDS':
        if move == 'LEFT':
            maze[x-1][y] = 'X'
        elif move == 'RIGHT':
            maze[x+1][y] = 'X'
        elif move == 'UP':
            maze[x][y-1] = 'X'
        elif move == 'DOWN':
            maze[x][y+1] = 'X'
    elif result == 'SUCCESS':
        maze[x][y] = 'X'
        print("Moved " + move + " from " + "(" + str(x) + "," + str(y) + ")")
        if backtracked == False:
            if move == 'LEFT':
                moves.append('RIGHT')
            elif move == 'RIGHT':
                moves.append('LEFT')
            elif move == 'UP':
                moves.append('DOWN')
            elif move == 'DOWN':
                moves.append('UP')