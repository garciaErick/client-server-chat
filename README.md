# client-server-chat
This is a sample multi-client/server done in golang. Has simple user logins and useranme validation. 

## How to run 
To run you simply need two terminals, one to act as a client, and the other one as server

### Server
You first need to start the server by either:
* Running the demoServer.sh script
    ```console
    me@garciaErick:/path/to/repo/$ ./demoServer.sh
    ```
* Running the following command
    ```console
    me@garciaErick:/path/to/repo/$ go run server.go authenticatedUser.go
    ```

### Client
Afterwards you can connect multiple clients to the server by
running the following in a separate terminal:
* Running the demoClient.sh script
    ```console
    me@garciaErick:/path/to/repo/$ ./demoClient.sh
    ```
* Running the following command
    ```console
    me@garciaErick:/path/to/repo/$ go run chatroom.go client.go  uuidGenerator.go
    ```
