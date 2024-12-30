# Go Lang CRUD API 🚀
Create Table Human

## Go lang <img src="./img/golang.png" width=30 height=30>
### Init Project Go Lang 🧑‍💻
```bash
go mod init api
```
### Install Pagkage 📥
```bash
go get github.com/gorilla/mux
go get github.com/lib/pq
```

## Run Docker Compose 🐳 
Go Lang(App) + Postgres(DB) + Nginx(Webserver)

### Run Docker Compose 💨
```bash
docker compose up -d 
```
### Up Scale Container Go-App 📈
```bash
docker compose up --scale go-app=3 --build
```

## How To Use API CRUD 📃

### Create 🔨
POST : ```localhost:8000/humans```

Body Raw
```
{
  "F_name":"frist Name"  
  "L_name":"Last Name"  
}
```
### Read 📖
all human\
GET : ```localhost:8000/humans```

select human {id}\
GET : ```localhost:8000/humans/{id}```

### Update 📝
PUT : ```localhost:8000/humans/{id}```

Body Raw
```
{
    "id":{id}
    "F_name":"frist Name"  
    "L_name":"Last Name"  
}
```

### Delete 💥
DELETE : ```localhost:8000/humans/{id}```
