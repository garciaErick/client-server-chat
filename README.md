# Go Client/Server Chatroom
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
## Demo
* The server supports multiple clients, and each connection is
aware of the current user, so messages from server are
personalized for each client.
![alt text](https://raw.githubusercontent.com/garciaErick/client-server-chat/master/screenshots/demo.PNG "Chatroom Demo")

* Additionally there is some simple validation, no duplicate
usernames are allowed and no illegal characters are allowed,
as an example ':' (colon) is an illegal character for a
username
![alt text](https://raw.githubusercontent.com/garciaErick/client-server-chat/master/screenshots/validation.PNG "Simple Validation Demo")

