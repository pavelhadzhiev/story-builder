# Story Builder

A Golang client-server application specifically used to play the story builder chat room game.

# Get Started

## Installation

Installing the Story Builder application is simple. As long as you have Golang installed, all you need to do is execute `go get github.com/pavelhadzhiev/story-builder`. Then, you will have the `story-builder` CLI (Command Line Interface) available. 

To display available commands and documentation, you can execute `story-builder help`.

To display documentation for a specific command you can execute `story-builder help <command>`.

## Connectivity

In order to play the story builder chat room game, you need to join a server. You can either connect to an existing one or host one on your own.

#### Host a Server

To host a server, a database is required. Currently the application uses strictly a [mysql](https://www.mysql.com/) database. Create one and execute `story-builder host <port> -u <dbUsername> -p <dbPassword>` where __port__ is the port at which you want to host the server and __dbUsername__ and __dbPassword__ provide the credentials for the database user. This command will start a server in the process that from which it's called. To kill it, press __ENTER__.

#### Connect to a Server

To connect to a server, execute `story-builder connect <hostname>` where __hostname__ is the host of the story builder server. This will check whether that server is online via a healthcheck endpoint, and then configure the CLI to use this server from now on.

#### Disconnect from a Server

To disconnect, execute `story-builder disconnect`. This will remove the server from the CLI configuration.

## Authentication

Once connected to a server, you need to either log in with an existing for the server user, or register a new one.

#### Log in

To log in, execute `story-builder login`. The command will prompt you to enter your username and password to log in with. It is also supported to pass credentials in flags, e.g. `login -u <username> -p <password>` if you need to log in from a script for example. The login command will check whether such a user exists in the server's user base and will configure the CLI to authenticate using this user from now on.

#### Register a New User

To register a new user, execute `story-builder register`. The command will prompt you to enter your username and password to register with. It is also supported to pass credentials in flags, e.g. `story-builder register -u <username> -p <password>` if you need to register a new user from a script for example. The register command will check whether such a user already exists in the server's user base, and if not - create it and configure the CLI to authenticate using this user from now on.

#### Log out

To log out, execute `story-builder logout`. This will erase the authentication from the CLI configuration.

## Game Rooms

Once you are successfully connected and authenticated, what's left is to join a game room and start playing.

#### List Rooms

To list the available rooms, execute `story-builder list-rooms`.

#### Join a Room

To join an existing room, execute `story-builder join-room <room>` where __room__ is the name of the room you want to join. This will set this room in the CLI configuration to allow gameplay commands.

#### Leave a Room

To leave your corrent room, execute `story-builder leave-room`. This will erase the current room from the CLI configuration and disable gameplay commands.

#### Create a Room

To create a new room, execute `story-builder create-room`. Initially you will be the only admin in your newly created room.

#### Delete a Room

To delete a room, execute `story-builder delete-room`. Note that this can only be done by the creator of the room.

#### Promote an Admin

In case you want to delegate someone else the access to manage your room, you can execute `story-builder promote <player>`. This command will give __player__ admin access in the current room. Note that this command itself requires admin access to be executed.

#### Ban a Player

In case you want to prevent someone from ever joining your room, you can execute `story-builder ban <player>` to ban the provided __player__. This requires admin access to be executed. If in the room, __player__ will be instantly remove and prevented from joining again.

## Gameplay

Once you are successfully connected, authenticated and you've joined a game room, a room admin may start a game. The server will let the users take turns and build up the story using the gameplay commands. If you need to check out your configuration, you can do so using the `story-builder info` command.

#### Start a Game

If you are an admin in the room, you can start a round of the game with `story-builder start-game`. This will give you the first turn so you can set the beginning, and create an order of the players to play after you. Once all players have played a turn, the order repeats until the game ends.

The `start-game` command can be executed with the `-e` or `--entries` flag and specify the number of turns the round will go on for. After they are played out the game will end automatically. If not used, the game goes on until the `end-game` command is executed by an admin.

The `start-game` command can be executed with the `-t` or `--time` flag to specify the time (in seconds) limit for a turn. If not used, the default value is 60 seconds.
 
The `start-game` command can be executed with the `-l` or `--length` flag to specifyr the maximum number of symbols that are allowed per entry. If not used, the default value is 100 symbols.

#### End a Game

If you are an admin in the room, you can end the game by executing `story-builder end-game`. This will notify users in the room that the next turn will be the last. After it is played the game ends.

Ending the game also support providing a specific number of turns to end after, e.g. `story-builder end-game 6`. This will end the game after 6 turns, instead of the next one.

#### Trigger a Vote Kick

If there are any problems with a specific player in the game, you can trigger a democratic vote process to kick him by executing `story-builder trigger-vote <player>` where __player__ is the player you want to kick. Once the vote treshold of __2/3__ of the players in the game is met, __player__ will be instantly kicked from the game. If it is his turn, it will be skipped to the next player.

#### Vote for an Ongoing Vote Process

If you want to cast your vote for the currently going voting, execute `story-builder vote`. Voting is allowed only once per player.

#### Get the Game

To get the status of the game, execute `story-builder get-game`. This will print all entries in the round so far in the order they were submitted and also show which users provided which entry. It will also show whose turn it is to play next, how much time is left for them to pay their turn, and how many more turns are to be played until the end of the game, if it was started with the `--entries` flag of the start command or `eng-game` was used to set it. It will also print information for any ongoing votes.

#### Add Entry to the Game

To add an entry to a story, execute `story-builder add <entry>` where entry is the text you wish to add to continue the story. Note that this requires that it is your turn and that your entry satisfies any game requirements (max quantity of symbols, etc.)

### Disclaimer

This project is part of the exam of a selective course in my university and is done with the sole purpose of learning and practicing Golang.
