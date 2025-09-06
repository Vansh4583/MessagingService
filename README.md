# Real-Time Messaging Service

A **real-time chat application** built in Go using a custom RPC framework over UDP.  
The system supports **user authentication**, **permission controls**, and **dual-mode message delivery**:  
- **Pull-based messaging** (clients poll for new messages manually)  
- **Push-based messaging** (server initiates callbacks to clients using remote object proxies)

The backend is organized into **independent microservices**:
- **Database Service (`dbd`)** – simple key-value store for service discovery
- **Authentication Service (`authd`)** – user signup/login, issues capabilities
- **Message Service (`messaged`)** – manages user inboxes, sends/receives messages, and pushes notifications
- **Chat Client (`chat`)** – command-line interface for users

---

## Features
- User signup and login with capability-based authentication
- Grant/block permissions for sending messages
- Pull messaging via `read` command
- Push messaging via `notify` command (server-initiated notifications)
- Service discovery through database registration

---



## How to Run

Each service runs as a separate Go process. **All processes pick a random port**, which gets registered with the database server for discovery.

### 1. Start the Database Server

```bash
go run cmd/dbd/dbd.go
```

- Outputs the database server's listening address (random port).  
- Example: `listening on 127.0.0.1:54072`

![DB Server](image_stub.png)

---

### 2. Start the Authentication Server

```bash
go run cmd/authd/authd.go <db_address>
```


- `<db_address>` = address printed by `dbd`.  
- Example: `go run cmd/authd/authd.go 127.0.0.1:54072`



![Auth Server](image_stub.png)

---

### 3. Start the Message Server

```bash
go run cmd/messaged/messaged.go <db_address>
```

- `<db_address>` = same as DB server above.  
- Example:  `go run cmd/messaged/messaged.go 127.0.0.1:54072`


![Message Server](image_stub.png)

---

### 4. Start Chat Clients
Open **two terminals** (for User A and User B). Each runs:

```bash
go run cmd/chat/chat.go <db_address> s <username> <password>
```


- `s` → signup and login (use `l` for login if already signed up).  
- Example:

```
go run cmd/chat/chat.go 127.0.0.1:54072 s alice pswd
go run cmd/chat/chat.go 127.0.0.1:54072 s bob pswd
```


![Chat Clients](image_stub.png)

---

## Commands Inside the Chat Client

```
s <user> <message> # Send message
read # Read all queued messages (pull mode)
notify # Enable push notifications (server will call back automatically)
allow <user> # Allow a user to send messages to you
block <user> # Block a user from sending messages
```



## Example Session
1. Alice allows Bob:

`allow bob`


2. Bob allows Alice:

`allow alice`


3. Alice enables push notifications:

`notify`


4. Bob sends Alice a message:

`s alice Hello!`


5. Alice immediately sees the pushed notification:

`bob: Hello!`



![Push Example](image_stub.png)

