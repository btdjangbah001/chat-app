# BASIC CHAT APP

## Functional requirements

### MVP

1. Users should be able to chat privately with other users
2. Users should be able to create chat groups and add other users by their usernames
3. Users should be able to have profile pictures
4. Users should be able to leave groups
5. Group owners must assign a successor before they can leave a group they own

### Extra Features

1. Users should be able to send media
2. There should be end to end encryption
3. Users should be able to post stories

## Building from source

1. Make sure you have go installed. You can get it from [Go offical site](https://go.dev/dl/)
2. Make sure you have Postgres DB installed. You can get it from [Postgres official site](https://www.postgresql.org/download/)
3. Clone the latest changes to your local machine

```
git clone https://github.com/btdjangbah001/chat-app.git
```

4. Run

```
go get
```

to install dependencies

5. Navigate to `/models/setup.go` and change database credentials by changing `dns` variable to

```
"host=<your_database_domain> password=<your_database_password> dbname=<your_database_name> port=<your_database_port> sslmode=disable"
```

6. To use a different database you can just change the database driver into whatever you want but note it was built with RDBMS in mind.

7. Run

```
go run .
```

to start the application.

8. If you are running it on your local machine you can access the app on

```
http://localhost:8080
```

### Pro-Tip

To test the app, there's a simple UI in `/chat-app-fe/index.html` you can use.
