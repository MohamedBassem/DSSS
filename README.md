# DSSS
Decentralized Secure Storage System


## System Components

### Upload Client
- Reads the Public/Private key pair for encryption.
- Divide, encrypt and upload files.
- Generate and upload the Manifest file.
- Asks the server "Where to upload" queries.

### Download Client
- Reads the Public/Private key pair for encryption.
- Fetches the manifest file from its first hash.
- Asks the server "Who Has" queries.
- Download file parts from the Manifest.
- Join parts and decrypt each part.

### Discovery Server
- Respond and delegates the "Who Has" queries from the client.
- Respond to "Where to upload" queries from the client.
- Respond to "IntroduceMe" queries from the client.
- Persistent connection with the agents.

### Agent
- Persistent connection with the discovery server.
- Respond to "Who Has" from the server.
- Respond to "InroductionRequests" from the server.
- Accept upload/get connections from peers.
- Store the hash/value pair on the FS.


## Protocol Definition

- Server/Client communication (HTTP Endpoints)
  - Who Has (GET /api/who-has?q=&lt;hash&gt;)
    - JSON Response `{ addresses : ["X.X.X.X:Z", "Y.Y.Y.Y:Q"] }`
  - Where to upload (GET /api/where-to-upload?size=&lt;size&gt;)
    - The response contain all the replicas to which the client should upload
    - JSON Response `{ addresses : ["X.X.X.X:Z", "Y.Y.Y.Y:Q"] }`
  - Introduce Me (GET /api/introduce-me?to=&lt;address&gt;&hash=&lt;hash&gt;&size=&lt;size&gt;)
    - The response contains the key the client should use to contact the agent
    - JSON Response `{ introduction-key : "X.X.X.X:Z" }`

- Server/Agent PING UPD connection for hole punching.

- Server/Agent communication (Persistent TCP connection)
  - Who has ("WHO_HAS &lt;hash&gt;")
    - Response "(0|1)"
  - Introduction Requests ("INTRODUCTION_REQUEST &lt;address&gt; &lt;size&gt; &lt;hash&gt;")
    - Response ("&lt;Introduction Key&gt;")

- Server/Agent communication (TCP over UDP using hole punching)
  - Upload Requests ("UPLOAD &lt;Introduction Key&gt; &lt;hash&gt; &lt;data&gt;")
    - Response ("Ok")
  - Download Requests ("DOWNLOAD &lt;hash&gt;")
    - Response ("&lt;data&gt;")

## Upload Flow

- Divide the file into X KBs parts.
- Encrypt each part and get its hash.
- For Each Part Do :
  - Store the hash into a manifest file.
  - Where to upload request to the discovery service.
  - For each address in the query response Do :
    - Introduction request.
    - Use the key in the response to send an upload request to the peer.
- Dump the manifest File, later ( TODO: Upload the file to the network )


## Download Flow

- Read the manifest file ( TODO: Fetch it from the network )
- For Each Line in the manifest Do :
  - Send a who has request.
  - Contact one of the peers, the server sent, with a download request . ( TODO: Authenticate first )
  - Decrypt file.
- Concat all parts and dump it.
