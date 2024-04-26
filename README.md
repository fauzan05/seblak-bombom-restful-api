# Aplikasi Warung Seblak Bombom

# Background

Starting from observing my friend's business, I then attempted to create an application to manage his business using the Go programming language.

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

# Database Table Structure

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
