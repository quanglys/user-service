#This is document of User-Service
##0. Project structure
- config: all configurations of this service
- doc: all documentations of this service
- requests: sample http requests to call this service's APIs
- run: contain executable file and configuration files, that can used to run
- scripts: all scripts like initial DB
- src: source code of this service

##1. APIs
- GetUser
- PostUser
- PatchUser
- GetUsers

 Read `api.yaml` for more detail about APIs

##2. Build project
- language:
```
required go v1.13 (can download and install via this site: https://golang.org/dl) 
```
- move to source project (user-service): 
```
source_proj=<source_proj>
cd ${source_proj}
```
- import library
```
go get github.com/go-kit/kit
go get github.com/gorilla/mux
go get github.com/spf13/viper
go get github.com/jinzhu/gorm
go get github.com/go-sql-driver/mysql
go get github.com/DATA-DOG/go-sqlmock
go get gotest.tools
```
- build
```
go build -o user-service ./src/
```

##3. Run service

- copy to destination dir
```
dest_dir=<dest_dir>
cp ./user-service ${dest_dir}
cp -r ./config ${dest_dir}
cd ${dest_dir}
```
- config: file ./config/config.yaml
```
http_server.port: port to bind service
mysql.uri: connection string is used to connect to mysql-db
```

- init mysql-db: 
```
run script in ./scripts/db-script.sql
```

- start service:
```
./user-service
```

- sample request: 
```
look and feel: ${source_proj}/requests/user-service.http 
```

##4. Todo

Other businesses and technical features should be implemented:

- More information about user like: address, date of birth, email, phone,...
- Manage relationships among users
- Manage permissions of users, check user's permission when calling service's APIs
- Add metrics to monitor service
- Integrate zipkin for better tracing if we get errors

If I have more time, I will do all above tasks 