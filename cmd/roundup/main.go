package main

import (
	"os"
	"fmt"
	"github.com/tomwright/monzoroundup/worker"
	"github.com/tomwright/monzoroundup/transport/http"
	"os/signal"
	"sync"
	"log"
	"github.com/tomwright/monzoroundup/monzo"
	"github.com/tomwright/monzoroundup/transport/http/auth"
	"github.com/tomwright/monzoroundup/user"
	"database/sql"
	"github.com/tomwright/monzoroundup/token"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(wg *sync.WaitGroup) {
		for range c {
			wg.Done()
		}
	}(&wg)

	monzo.Init(fmt.Sprintf("%s://%s%s", os.Getenv("HTTP_PROTOCOL"), os.Getenv("PUBLIC_DOMAIN"), auth.CallbackEndpoint))

	db := getDB()
	userModel := user.NewSqlModel(db)
	tokenModel := token.NewSqlModel(db)

	go worker.Listen(tokenModel, userModel)

	http.Init(os.Getenv("HTTP_BIND_ADDRESS"), tokenModel, userModel)
	go http.Listen()

	wg.Wait()
	log.Println("Exiting...")
	os.Exit(0)
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	return db
}
