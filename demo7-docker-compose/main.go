package main

import (
	"context"
	"demo7-docker-compose/endpoint"
	"demo7-docker-compose/model"
	"demo7-docker-compose/service"
	"demo7-docker-compose/transport"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	servicePort := flag.Int("service.port", 8088, "service port")

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	model.ConnMysql("10.10.20.16", "3306", "root", "123456", "user")
	model.ConnRedis("10.10.20.16", "6389", "")

	userService := service.MakeUserServiceImpl(&model.UserDaoImpl{})

	r := transport.MakeHttpHand(ctx, &endpoint.UserEndpoints{
		RegisterEndpoint: endpoint.MakeRegisterEndpoint(userService),
		LoginEndpoint:    endpoint.MakeLoginEndpoint(userService),
	})

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err := <-errChan
	fmt.Println(err)
}
