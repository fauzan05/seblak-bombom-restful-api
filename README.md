# Seblak Bombom Eatery App

# Background

Starting from observing my friend's business, I then attempted to create an application to manage his business using the Go programming language. This is a "seblak" sales application where payments can be made onsite or online. For onsite payments, customers only need to come to the location and an admin will input their data manually to indicate whether they have paid or not. Meanwhile, for online payments, customers simply choose from the payment methods provided by Xendit, and the payment status will be automatically displayed. So the admin no longer needs to manually input data because everything is handled by Xendit. For delivery distance, we use kilometers as the unit, and there is a menu for delivery service rates along with the delivery status. This application is integrated with the Xendit payment gateway. Make sure you have read the documentation from Xendit to understand the payment flow further.

<br>

# Requirements

## Tech Stack

- Go version 1.22.0  : https://go.dev/
- MySQL version 5.7 : https://mariadb.org/

## Frameworks & Libraries

- Fiber version 2.0 : https://gofiber.io/
- GORM version 2.0 : https://gorm.io/index.html
- Validator version 10 : https://pkg.go.dev/github.com/go-playground/validator/v10
- Migrate version 4.17.0 : https://github.com/golang-migrate/migrate
- Viper : https://github.com/spf13/viper
- Logrus : https://github.com/sirupsen/logrus
- Testify : https://github.com/stretchr/testify
- MySQL : https://github.com/go-sql-driver/mysql
- Xendit : https://github.com/xendit/xendit-go

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

After you typing ``` migrate ```, it will shown Database drivers : stub, mysql like below :
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
Before proceeding further, make sure you have an account on Xendit. Just create one, and it's free. After finishing creating the account, input the environment such as **Secret Key**, **Public Key**, **Business Id** and **Callback Token** in the config.json file located in the root directory. Because config.json is listed on .gitignore file, so you have to create another one with copying the config-example.json and then renamed it into config.json. So, after that because I'm using Docker Compose environment, first you have to create a docker network so this app can communicate with api consumer (seblak-bombom-api-consumer) inside the docker environment. What needs to be done is to execute the following command
```
docker network create seblak-bombom-network
```
```
docker compose build
```
This command will build the images specified in the Dockerfile and MariaDB. Then, once it's done, execute the following command:
```
docker compose up
```
This command is used to run the configurations that have been set up in the Dockerfile as well as the MariaDB container. At this stage, when we make changes to any files in this application, they will be automatically recompiled without the need for manual compilation. This is because I'm using CompileDaemon, which has been fetched in the Dockerfile.
<br>
docker-compose logs -f


If you're not using Docker, simply run **go run main.go** in the app directory. Beforehand, make sure the database is running and the database is created according to the configuration in the config.json file in the root directory. Once done, perform operations on the API by referring to api-specs.json to understand the request and response of each endpoint.
<br>

For testing purposes on the Midtrans endpoint, I recommend using Ngrok because Xendit Callback Notification, after the user completes the payment, requires the redirection settings to be an active URL (endpoint) accessible over the internet. For example, by exposing a URL like this and entering it in the finish URL section:

**https://6901-180-243-9-232.ngrok-free.app/api/xendits/payment-request/notifications/callback**

So after a successful transaction or expiration, Xendit will automatically perform a callback to that endpoint.

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

# Important

if you want to login with admin, you must create an admin account with the path "/users/register", but you must include "Custom Header API Key Authentication" with the name "X-Admin-Key" then followed by the request body in the role = admin section. For X-Admin-Key, you must set it in your env file in the "ADMIN_CREATION_KEY" section and fill in whatever you think is secret. For example, it is in the .env.example file and you copy it and then name it .env .

## 🧪 Live Demo / Testing

You can try this application directly on the following demo server:

🔗 **Frontend**: [https://seblak.fznh-dev.my.id](https://seblak.fznh-dev.my.id)
<br>
🔗 **Backend (API base URL)**: [https://api.fznh-dev.my.id](https://api.fznh-dev.my.id/api)

### 🔐 Demo Login (User/Customer)
- Email: `cust1@email.com`
- Password: `Cust1Testing#`

### 🔐 Demo Login (Admin)
- Email: `admin1@email.com`
- Password: `Admin1Testing#`

### 🚧 Note
- Data on the demo server will be reset periodically
- You **don't need to register manual**, just use the available demo account
- Some features (such as payments) are simulated/mock

