# Seblak Bombom Eatery App

# Background

Starting from observing my friend's business, I then attempted to create an application to manage his business using the Go programming language. This is a "seblak" sales application where payments can be made onsite or online. For onsite payments, customers only need to come to the location and an admin will input their data manually to indicate whether they have paid or not. Meanwhile, for online payments, customers simply choose from the payment methods provided by Midtrans, and the payment status will be automatically displayed. So the admin no longer needs to manually input data because everything is handled by Midtrans. For delivery distance, we use kilometers as the unit, and there is a menu for delivery service rates along with the delivery status. This application is integrated with the Midtrans payment gateway. Make sure you have read the documentation from Midtrans to understand the payment flow further.

<br>

# Requirements

## Tech Stack

- Go version 1.22.0  : https://go.dev/
- MariaDB version 10.4.28 : https://mariadb.org/

## Frameworks & Libraries

- Fiber version 2.0 : https://gofiber.io/
- GORM version 2.0 : https://gorm.io/index.html
- Validator version 10 : https://pkg.go.dev/github.com/go-playground/validator/v10
- Migrate version 4.17.0 : https://github.com/golang-migrate/migrate
- Viper : https://github.com/spf13/viper
- Logrus : https://github.com/sirupsen/logrus
- Testify : https://github.com/stretchr/testify
- MySQL : https://github.com/go-sql-driver/mysql
- Midtrans : https://midtrans.com/

## Software
- Docker
- Terminal
- Postman
- Ngrok

<br>

# Installation
Please install/add packages in the root directory project first below :

Install Framework Go Fiber

```
go get github.com/gofiber/fiber/v2
```

Install GORM

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

Install Validator

```
go get github.com/go-playground/validator/v10
```

Install Go Migrate

```
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
<!-- type 'migrate' without '' to check if migration mysql has been installed -->
```
If mysql in migration has been installed, you must have to set GOPATH so it can be call the command with 'migrate' in terminal. You can check this post to understand how to add GOPATH :
https://stackoverflow.com/questions/21499337/cannot-set-gopath-on-mac-osx

After you typing 'migrate' without '', it will shown Database drivers : stub, mysql like below :
![alt text](<Screenshot 2024-03-26 at 20.38.04.png>)


Install Viper

```
go get github.com/spf13/viper
```

Install Logrus

```
go get github.com/sirupsen/logrus
```

Install Testify

```
go get github.com/stretchr/testify
```

Install Mysql Driver

```
go get -u github.com/go-sql-driver/mysql
```

Install UUID

```
go get github.com/google/uuid
```

<br>

# How to run
Before proceeding further, make sure you have an account on Midtrans. Just create one, and it's free. After finishing creating the account, input the environment such as **MerchantID**, **Client Key**, and **Server Key** in the config.json file located in the root directory. So, after that because I'm using Docker Compose environment, what needs to be done is to execute the following command:
```
docker compose build
```
This command will build the images specified in the Dockerfile and MariaDB. Then, once it's done, execute the following command:
```
docker compose up
```
This command is used to run the configurations that have been set up in the Dockerfile as well as the MariaDB container. At this stage, when we make changes to any files in this application, they will be automatically recompiled without the need for manual compilation. This is because I'm using CompileDaemon, which has been fetched in the Dockerfile.
<br>

If you're not using Docker, simply run **go run main.go** in the app directory. Beforehand, make sure the database is running and the database is created according to the configuration in the config.json file in the root directory. Once done, perform operations on the API by referring to api-specs.json to understand the request and response of each endpoint.
<br>

For testing purposes on the Midtrans endpoint, I recommend using Ngrok because Midtrans Callback Notification, after the user completes the payment, requires the redirection settings to be an active URL (endpoint) accessible over the internet. For example, by exposing a URL like this and entering it in the finish URL section:

**https://8xx2-xxx-xx2-1xx5-9x1.ngrok-free.app/api/midtrans/snap/orders/notification**

So after a successful transaction or expiration, Midtrans will automatically perform a callback to that endpoint.

# Go Migrate Command

To make a migration :
```
migrate create -ext sql -dir database/migrations create_table_xxx
```

To run a migrations :
```
// run migration
migrate -database "mysql://root@tcp(localhost:3306)/database_name" -path database/migrations up

// rollback migration
migrate -database "mysql://root@tcp(localhost:3306)/database_name" -path database/migrations down
```

To remove dirty :
```
// V is a version of dirty column
migrate -path database/migrations -database "mysql://root@tcp(localhost:3306)/database_name" force V
```

To migrate 1 step :

```
// 1 is a how many step do you wanna
migrate -database "mysql://root@tcp(localhost:3306)/database_name" -path database/migrations up 1
```
