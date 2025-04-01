# WASAText
This is a project for the Web and Software Architecture course from Sapienza University. It's a simple messaging app. The project defines APIs using the OpenAPI standard, contains the server side - backend in Go (uses SQLite database), the frontend in JavaScript and creates a Docker container image for deployment.
The main features of the app are the following:

-Conversations List: Users see chats sorted in reverse chronological order, displaying the other user's name or group name, latest message preview, and timestamp.

-Messaging: Users can send, reply to, forward, and delete messages. Messages include timestamps, sender info, and a status indicator (✓ received). Reactions (emoticons) can be added or removed.

-Groups: Users can create groups, add members, and leave at any time. Only existing members can invite others.

-User Search & Management: Users can search for existing WASAText usernames.

-User profile: The user can enter their profile page, update their username and profile photo.

-Authentication: Users log in with just a username—if it exists, they are logged in; otherwise, they are registered automatically. The API returns a user identifier, which is used as a token in the Authorization header - the project uses bearer authentification

## To run backend only

Navigate to the repository directory, then run:
```shell
go run ./cmd/webapi/
```
After running the backend functionalities can be tested by sending HTTP requests (using the terminal) to the server running on localhost at port 3000. For example this request instructs the backend to add an user of a specified username to an existing group.

```shell
curl -X POST http://localhost:3000/group/1/add \
-H "Authorization: Bearer 4" \
-H "Content-Type: application/json" \
-d '{"username": "giulia"}'
}'
```
## Database
In order to check the contents of the database, run these commands:
```shell
Cd /tmp
```
```shell
Sqlite3 decaf.db 
```
Then it is possible to run sql queries, for example to check the users table:
```shell
SELECT * FROM users;
```

## To run the WebUI (for production)

```shell
./open-node.sh

yarn run build-prod

yarn run preview
```

## To run with Docker

```shell
docker build -t wasa-text-frontend:latest -f Dockerfile.frontend .
docker build -t wasa-text-backend:latest -f Dockerfile.backend .
```
```shell
docker run -it --rm -p 3000:3000 wasa-text-backend:latest 
docker run -it --rm -p 8080:80 wasa-text-frontend:latest
```

## Future improvements

-Fixing the profile/group photo display functionality

-Implementing double checkmarks

-When a message is forwarded, the message preview doesn’t update to the forwarded message

## Error checks

How to check for errors like in the homework evaluations:

1.	Open-api linter:
   
https://github.com/IBM/openapi-validator

Enter the directory and run on api.yaml file as instructed in above repository.

The open-api linter checks the file against yaml and openapi specifications.

2.	Go linter:
   
Golangci-lint

run golangci-lint run service/database --enable-all –verbose

run golangci-lint run service/api --enable-all –verbose

With this setting you can find all the errors and corresponding lines of code in which the errors arise.

To check rowserr:

golangci-lint run service/api -E rowserrcheck

golangci-lint run service/database -E rowserrcheck

3.	Vuejs:

Eslint vuejs linter: https://eslint.vuejs.org/

Although I would recommend to have a backup of the latest version of the project or run this on a cloned virtual machine since using this created some issues with the code for me

Use Prettier for code formatting
