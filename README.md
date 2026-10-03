# Chirpy!
Chirpy is a web server that simulates twitter! It has all the basic functions like registering/login/authenticating users and publishing and deleting chirps (the equivalent of a chirp).

## Instalation
### Go
The chirpy server is a Go app so you'll need the go toolchain to compile it.
One pretty easy way to get the go toolchain is downloading it from webi, simply run `curl -sS https://webi.sh/golang | sh; source ~/.config/envman/PATH.env`

 You can copy the repo and compile the code using the next command on the bash:

```shell
git clone https://github.com/AyuorusAguilar/chirpy.git && go build . -o chirpy_server
```

This will leave an executable for the chirpy server on your working directory. After executing it, it will start listening at the 8080 port, but before doing that you need to set up some things, like the database and some variables.

### Database
You can set up your own database for chirpy by following the next steps
- Download Postgres. You can install postgres from the default package manager for Debian (apt) by running: `sudo apt install postgresql postgresql-contrib`
- Create a new database for chirpy. Connect using `sudo -u postgres psql` and there run `CREATE DATABASE chirpy;`
- You'll need a password so create a user or just set a password for the default user postgres using the command `ALTER USER postgres PASSWORD 'postgres';`.
- Get your connection sting! Just replace the fields with your username (or with "postgres" if you're using the default) and password `postgres://<username>:<password>@localhost:5432/chirpy`
- Write your connection string inside the `.chirpy_conf.json` file at the root of the project
- Finally, to set up the database you can use goose and run the migrations in the `sql/schema/` folder of this project. Simply install it using the go toolchain with `go install github.com/pressly/goose/v3/cmd/goose@latest` and then navigate to the `sql/schema/` directory and run `goose postgres <connection_string> up` just replace the connection string. Alternatively you could just manually run the migration's queries in the `sql/schema/` folder, but I think that would be pretty annoying!.

After doing all that, you can just run the executable and the server will be working!

## Usage
### Admin
There is only one endpoint tought to be used by admins of the server:
- **POST /admin/reset**: This method resets the database, so use it carefully!

### User
#### Authentication
User interact with the server with the following commands, listed in the order they'll probably follow:
- **POST /api/users**: This method registers a user in the database.
```json
{"email": "example@exmp.com", "password":"totallysecurepasswordofcourse69"}
```

- **POST /api/login**: This method registers a login for the given user and returns both a refresh_token and a JWT token. The refresh token represents a login lasts 60 days, you can revoke it using the `POST /api/revoke` endpoint, until you do that or until it expires, any request that uses your refresh_token will be accepted and treated as coming from that user. The JWT Token, on the other side, only last one hour and is the one used to access other parts of the server, you can get another one using the `POST /api/refresh` and your refresh token. This *login* endpoint, by the way, expects a body like this one:
```json
{"email": "example@exmp.com", "password":"totallysecurepasswordofcourse69"}
```

- **POST /api/refresh**: This method gives you a new JWT token, but it requieres an *Authorization* Header like `Bearer <refresh_token>`, you can get a refresh token using the `POST /api/login` endpoint.

- **POST /api/revoke**: You can log out calling this endpoint with an *Authorization* Header like `Bearer <refresh_token>`

- **PUT  /api/users**: This method updates a user's cretentials. It expects an *Authorization* Header like `Bearer <token>` body like this one:
```json
{"email": "example@exmp.com", "password":"totallysecurepasswordofcourse69"}
```
#### Chirping!

- **POST /api/chirps** : Post a chirp! All this endpoint expects is an *Authorization* Header like `Bearer <refresh_token>` and a body like this one:
```json
{"email": "example@exmp.com", "password":"totallysecurepasswordofcourse69"}
```

- **GET  /api/chirps/** : This endpoint returns all chirps in the database. You don't need an authorization token to hit this endpoint. You can specify some optional query parameters like "author_id=<>" that will return only the chrips made by that specific UUID or "sort=<asc/desc>" that will order the chirps by creation date.


- **GET  /api/chirps/{ChirpId}** : This endpoint will return the chirp object of the specified UUID

- **DELETE /api/chirps/{ChirpId}** : This endpoint will delete the chirp of the specified UUID as long as the Token in the *Authorization* Header like `Bearer <refresh_token>` matches the author of the chirp.

### Other
- **GET  /api/healthz**: It returns code 200 if the server is up and running