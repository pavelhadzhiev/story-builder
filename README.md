# Story Builder

A Golang client-server application specifically used to play the story builder chat room game.

# Get Started

## Installation

Installing the Story Builder application is simple. As long as you have Golang installed, all you need to do is execute `go get github.com/pavelhadzhiev/story-builder`. Then, you will have the `story-builder` CLI (Command Line Interface) available. 

To display available commands and documentation, you can execute `story-builder help`.

To display documentation for a specific command you can execute `story-builder help <command>`.

## Connectivity

In order to play the story builder chat room game, you need to join a server. You can either connect to an existing one or host one on your own.

#### Connect to a Server

To connect to a server, execute `story-builder connect <server>` where server is the host of the story builder server. This will check whether that server is online via a healthcheck endpoint, and then configure the CLI to use this server from now on.

#### Host a Server

To host a server, execute `story-builder host <port>` where port is the port at which you want to host the server. This will start a server in a separate process. To kill it, press __ENTER__.

#### Disconnect from a Server

To disconnect, execute `story-builder disconnect`. This will remove the server from the CLI configuration.

## Authentication

Once connected to a server, you need to either log in with an existing for the server user, or register a new one.

#### Log in

To log in, execute `story-builder login -u <username> -p <password>`. This will check whether such a user exists in the server's user base and will configure the CLI to authenticate using this user from now on.

#### Register a New User

To register a new user, execute `story-builder register -u <username> -p <password>`. This will check whether such a user already exists in the server's user base, and if not create it and configure the CLI to authenticate using this user from now on.

#### Log out

To log out, execute `story-builder logout`. This will erase the authentication from the CLI configuration.

## Game Rooms

Once you are successfully connected and authenticated, what's left is to join a game room and start playing.

#### List Rooms

To list the available rooms, execute `story-builder list-rooms`.

#### Join a Room

To join an existing room, execute `story-builder join-room <room>` where room is the name of the room you want to join. This will set this room in the CLI configuration to allow gameplay commands.

#### Leave a Room

To join an existing room, execute `story-builder leave-room`. This will erase the current room from the CLI configuration and disable gameplay commands.

#### Create a Room

To create a new room, execute `story-builder create-room`.

#### Delete a Room

To delete a room, execute `story-builder delete-room`. Note that this can only be done by the creator of the room.

## Gameplay

Once you are successfully connected, authenticated and you've joined a game room, the room owner may start a game. The server will let the users take turns and build up the story using the gameplay commands.

#### Start a Game

If you are the room owner you can start a round of the game with `story-builder start-game`. This will give you the first turn and create an order of the players to play after you. Once all players have played a turn, the order repeats until the game ends. Newly joined users will enter at the end of the order of play.

The `start-game` command can be executed with a `-t` or `--turns` flag and specify the number of turns the round will go on for. After they are played out the game will end automatically. If not used, the game ends when the `end-game` command is executed by the room owner.
 
#### End a Game

If you are the room owner, you can end the game by executing `story-builder end-game`. This will notify users in the room that the next turn will be the last. After it is played the game ends.

#### Print the Story

To get the status of the game, execute `story-builder story`. This will print all entries in the round so far in order and also show which users provided which entry. It will also show whose turn it is to play next and how many more turns are to be played until the end of the game, if it was started with the `--turns` flag of the start command.

#### Add to a Story

To add an entry to a story, execute `story-builder add <entry>` where entry is the text you wish to add to continue the story. Note that this requires that it is your turn and that your entry satisfies any game requirements (max quantity of symbols, etc.)

#### Get Users

To check available users, execute `story-builder users`. This will show you which users are currently in your game and whose turn it is.

### Disclaimer

This project is part of the exam of a selective course in my university and is done with the sole purpose of learning and practicing Golang.
