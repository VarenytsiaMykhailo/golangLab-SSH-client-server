package main

//https://godoc.org/github.com/gliderlabs/ssh
import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
)

func main() {
	ssh.Handle(func(s ssh.Session) { //передаваемая функция - обрабатывает установленные сеансы

		io.WriteString(s, "Hello world\n")
	})

	//Запуск сервера.
	log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
		return pass == "secret" //тут можно перечислить пароли, с которыми будет пускать клиента на сервер (например secret)
	})))

	/*
	2ой параметр - функция handler-обработчки. nil - юзать дефолтный
	***
	Если вы не укажете ключ хоста, он будет генерироваться каждый раз. Это удобно, за исключением того,
	что вам придется иметь дело с клиентами, которые путаются в том, что ключ хоста отличается.
	Лучше создать или указать существующий ключ в вашей системе:
	log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.HostKeyFile("/Users/progrium/.ssh/id_rsa")))
	*/

}
