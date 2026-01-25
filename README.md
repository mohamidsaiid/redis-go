# Redis-Go

This is a toy Redis server implementation in Go, built as part of the CodeCrafters "Build your own Redis" challenge.

## Description

This project is a simplified implementation of a Redis server that supports a subset of Redis commands. It's designed to be a learning exercise in building concurrent network applications and understanding the Redis protocol (RESP).

## Implemented Commands

The following Redis commands are supported:

### Connection Management

*   `PING`: Returns "PONG", used to test if the server is alive.
*   `ECHO message`: Returns the given message.

### String Operations

*   `SET key value [EX seconds | PX milliseconds]`: Sets the string value of a key. The optional `EX` and `PX` flags set an expiry time on the key.
*   `GET key`: Gets the string value of a key.

### List Operations

*   `LPUSH key element [element ...]`: Prepends one or more elements to a list.
*   `RPUSH key element [element ...]`: Appends one or more elements to a list.
*   `LPOP key [count]`: Removes and returns the first element of a list. An optional `count` can be provided to pop multiple elements.
*   `LRange key start stop`: Returns the specified elements of the list stored at the key.
*   `LLEN key`: Returns the length of the list stored at the key.
*   `BLPOP key timeout`: Is a blocking list pop operation. It blocks the connection when there are no elements to pop from the list and unblocks when new elements are pushed or when the timeout is reached.

### Introspection

*   `TYPE key`: Returns the type of the value stored at the key.

## Data Storage

The server uses a thread-safe, in-memory key-value store. This is achieved by using Go's `sync.Map`, which allows for safe concurrent access from multiple goroutines without the need for explicit locking.

### Performance and Concurrency

*   **Concurrency Model**: The server is highly concurrent. It spawns a new goroutine for each incoming client connection. Each command is also handled in its own goroutine, allowing the server to process multiple commands from multiple clients simultaneously.
*   **Data Structures**: The use of `sync.Map` is a key performance consideration. It is optimized for the case where keys are written once and read many times, which can be a common pattern in some Redis use cases. For list operations, the server uses Go's built-in slices. While this is a simple approach, for very large lists, a more optimized data structure like a linked list could provide better performance for push and pop operations at the ends of the list.
*   **`BLPOP` Implementation**: The implementation of `BLPOP` uses a channel (`newPushedElement`) to signal when a new element has been pushed to a list. This is a simple and effective way to implement blocking operations without busy-waiting. However, in a production system with many blocked clients, a more sophisticated pub/sub or notification system might be necessary to avoid contention on the channel.

## How to Run

The server can be run using the provided shell script:

```sh
./your_program.sh
```

This will start the Redis server on `0.0.0.0:6379`. You can then connect to it using `redis-cli`.

## How to Run Tests

The project is set up to be tested with the CodeCrafters testing platform. The tests can be run by pushing the code to the CodeCrafters git remote.

The `codecrafters.yml` file defines the test stages. The `your_program.sh` script is the entry point for running the server, and the tests will interact with the running server to verify its functionality.