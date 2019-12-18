package main

//https://godoc.org/github.com/gliderlabs/ssh
import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"io/ioutil"
	"log"
	"strconv"
)

var rootPath = "./server/root/"

func main() {
	sl := []string{"1", "2", "3"}
	fmt.Println(sl)
	ssh.Handle(func(s ssh.Session) { //передаваемая функция - обрабатывает установленные сеансы
		//defer s.Exit()
		sliceOfCommands := s.Command()
		if sliceOfCommands[0] == "ls" {
			dirsAndFiles, err := ioutil.ReadDir(rootPath) //инфа по содержимому в текущей папке (получаемый слайс - уже в отсортированном по имени виде)
			if err != nil {
				io.WriteString(s, "err") //отправляем ошибку клиенту
			}
			var dirs string  //сюда заносим названия папок в директории Path
			var files string //сюда заносим названия файлов в директории Path
			for _, file := range dirsAndFiles { //перебор содержимого текущей папки
				if file.IsDir() {
					dirs += file.Name() + "\n"
				} else { //если это файл, а не папка
					files += file.Name() + " (" + strconv.Itoa(int(file.Size())) + "b)\n"
				}
			}
			io.WriteString(s, dirs+files)
			return
		}
	})

	//Запуск сервера.
	log.Fatal(ssh.ListenAndServe("127.0.0.1:2222", nil, ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
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
