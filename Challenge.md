# Mazes
## UPE Fall 2018 Coding Challenge

### Objective
UCLA just defeated Cal football 37-7. However, you don't care because when you went to the career fair you were
really hitting it off with a super friendly recruiter for your dream job who was about to hand you a full-time
offer/internship with $1,337,420 starting salary! But right before they could hand you the offer to sign, HKN
and TBP officers dressed in USC colors jump the recruiter and kidnap them. 

Luckily, Professor Smallberg has come through for you many times and this time is no exception: he has discovered
a way to precisely map the path with which the HKN and TBP officers took as a maze. He has linked this maze directly 
with your brain, but due to some bad circuitry you are not able to see the maze. The only way you are able to navigate 
is use his API to send HTTP requests to tell the direction you intend to move, and Professor Smallberg will whisper back 
to you whether you successfully moved in that direction, if you encountered a wall, or encountered the edge of the maze.

### Technical Summary
To begin the challenge, you must make a POST call with your Student UID to http://ec2-34-216-8-43.us-west-2.compute.amazonaws.com/session. This call will return to you a token that you can use to make future calls to affect the game state. The game session will expire after 5 minutes; as a result, you must write a program to solve the maze, otherwise you will not be able to manually solve all the randomly created mazes in time before your token expires.

To get details on the maze state, you are then able to make a GET request with your token to http://ec2-34-216-8-43.us-west-2.compute.amazonaws.com/game?token=[ACCESS_TOKEN]. This will tell you the size of the maze, the current location that you're at, the status of the maze, how many levels you have completed, and the total number of levels. This state will change in real time as you take actions to move through the maze. To move through the maze, make POST requests to http://ec2-34-216-8-43.us-west-2.compute.amazonaws.com/game?token=[ACCESS_TOKEN] with a request body containing an action with 4 possible moves (LEFT, UP, RIGHT, DOWN).

#### Example Maze
<img width="650" alt="screen shot 2018-10-14 at 8 17 44 pm" src="https://user-images.githubusercontent.com/12661925/46928265-96903780-cfee-11e8-9fad-50f7ef21f8f6.png">

Internal Representation:
```
[
	“S**  *   “,
	“  * *  * “,
	“ ** ** * “,
	“     *E* “,
	“ *** *** “,
	“*      * “,
	“  **** * “,
	“ *   *** “,
	“   *     “
]
```


### API
The base url is http://ec2-34-216-8-43.us-west-2.compute.amazonaws.com

#### POST /session
```
Request Body:
{
	“uid”: str
}

Expected Response Body:
{
	“token”: str # token encoded with uid
}
```
- When you begin a session, you will be returned the same token until the session expires. 
- To start a new session, wait for the session to time out.

#### GET /game?token=[ACCESS_TOKEN]
```
Expected Response Body:
{
	“maze_size”: [int, int], <- [width, height], null if status is NONE or FINISHED
	“current_location”: [int, int], <- [xcol, ycol], null if status is NONE or FINISHED
	“status”: str, <- can be “PLAYING”, “GAME_OVER”, “NONE”, “FINISHED”
	“levels_completed”: int <- 0 indexed, 0-L, null if status is NONE or FINISHED
	“total_levels”: int <- L, null if status is NONE or FINISHED
}
```
- State will be NONE if your session has expired or does not exist
- A new game state is returned whenever a new session is created
- A new game state is returned whenever you complete a game/maze

#### POST /game?token=[ACCESS_TOKEN]
```
Request Body:
{
	“action”: str
}

Expected Response Body
{
	“result”: str
}
```
- The result is “WALL” if a wall is hit, “SUCCESS”, “OUT_OF_BOUNDS” if outside, or “END” if end has been reached
- Action must be “UP”, “DOWN”, “LEFT”, “RIGHT”


### Notes
- You may use any language for your solution.

- Standard HTTP codes and JSON messages will be returned for errors, e.g. invalid action in GET /game?token=str&action=[ACTION] will trigger 400 and a relevant error message.

- Please only use your own Student UID.

- Please do not hit Professor Smallberg's API endpoints too hard, as your mind activity may overload and you may get an aneurysm!

- To pass a game, you only need to find your way to the end of the maze from the start.

- Please ask questions in the comments section below so that we can respond and help everybody!

- Be sure to replace the access token in your requests with the new active token every time that your session expires!

- You will have the entire quarter until Friday of Week 8 to complete the task. Feel free to ask your Byte for help or guidance on the task -- they are here to help you grow as a developer and follow an agile development process -- but do not work with other inducting bits.

### Submitting Your Solution
When you are done, share a GitHub repository containing your solution code to **austinguo550** or **Austin Guo** on GitHub. Your submission time will be based on this email, though you may continue fixing up your code and making it more maintainable. Your completion status will be read from the server for grading at some arbitrary time after Friday of 9th week at 11:55:00 pm PST, so be sure to complete the challenge before that time.

Please also include sufficient documentation on how to build and run your submission in a README file in your repository.