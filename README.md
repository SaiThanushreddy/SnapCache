# Redis-like In-Memory Database

This project implements an in-memory database compatible with the Redis protocol using the RESP (REdis Serialization Protocol). The server listens on port 6379 and can handle basic commands.

## Features

- Basic RESP commands (e.g., PING)
- Connection handling
- Support for future data persistence with AOF (Append Only File)

## Prerequisites

- Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).

## Running the Server

1. Clone this repository:

   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Open a terminal and navigate to the project directory.

3. Run the server:

   ```bash
   go run main.go
   ```

   You should see the message: `Listening on port :6379`.

4. You can test the server using telnet or any Redis client. For example:

   ```bash
   telnet localhost 6379
   ```

   Once connected, you can type:

   ```
   PING
   ```

   You should receive a response:

   ```
   +PONG
   ```

## Currently Working On

### AOF (Append Only File)

This project aims to implement the AOF (Append Only File) method for data persistence. AOF records every command executed in the database in a RESP format. When the server starts, it reads the commands from the AOF file and executes them to restore the database state.

The implementation involves creating an Aof struct, which will hold the file and a buffered reader to read RESP commands. The NewAof function will initialize the AOF file and a goroutine will periodically sync the AOF file to ensure data durability.

#### Summary of AOF Implementation

Data persistence is critical for any database, including in-memory databases. The AOF method provides a way to keep records by logging each command executed.

##### What is AOF?
- AOF (Append Only File) records every command as it is executed in the database.
- It allows recovery of the database state by re-executing these commands after a crash or server restart.

##### AOF Format
Commands are stored in the AOF file in the RESP format. For example, if we execute the following commands:

```
SET name ahmed
SET website ahmedash95.github.io
```

The AOF file would contain:

```
*2
$3
SET
$4
name
$5
ahmed
*3
$3
SET
$7
website
$20
ahmedash95.github.io
```

##### Writing the AOF Struct
- Create an Aof struct to handle file operations.
- Implement the NewAof method to initialize the file and a buffered reader.
- Use a goroutine to periodically sync the AOF file to disk, ensuring durability even in the event of a crash.
