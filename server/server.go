package main

//https://godoc.org/github.com/gliderlabs/ssh
import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var rootPath = "./server/root/"

func main() {
	sl := []string{"1", "2", "3"}
	fmt.Println(sl)
	ssh.Handle(func(s ssh.Session) { //передаваемая функция - обрабатывает установленные сеансы
		//defer s.Exit()
		sliceOfCommands := parseCommands(s) //парсим команды клиента

		switch sliceOfCommands[0] {
		case "ls":
			listDir(s, sliceOfCommands)
		case "mkdir":
			mkdir(s, sliceOfCommands)
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

func parseCommands(s ssh.Session) []string {
	strOfCommands := s.RawCommand()
	firstIndexOfSpace := strings.Index(strOfCommands, " ")
	var sliceOfCommands []string //parsed
	if firstIndexOfSpace != -1 {
		sliceOfCommands = append(sliceOfCommands, strOfCommands[0:firstIndexOfSpace])
		sliceOfCommands = append(sliceOfCommands, strOfCommands[firstIndexOfSpace+1:])
	} else {
		sliceOfCommands = append(sliceOfCommands, strOfCommands[0:])
		sliceOfCommands = append(sliceOfCommands, "")
	}
	return sliceOfCommands
}

func listDir(s ssh.Session, sliceOfCommands []string) {
	var path = ""
	if sliceOfCommands[1] != "" {
		if sliceOfCommands[1][0:2] != "./" {
			io.WriteString(s, "incorrect path. Use ex: \"./test dir/asd.txt\"") //отправляем ошибку клиенту
			return
		}
		path = sliceOfCommands[1][2:]
	}
	dirsAndFiles, err := ioutil.ReadDir(rootPath + path) //инфа по содержимому в текущей папке (получаемый слайс - уже в отсортированном по имени виде)
	if err != nil {
		io.WriteString(s, err.Error()) //отправляем ошибку клиенту
		return
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
}

func mkdir(s ssh.Session, sliceOfCommands []string) {
	var path = ""
	if sliceOfCommands[1] != "" {
		path = sliceOfCommands[1]
	} else { //если название папки не было переданно
		io.WriteString(s, "incorrect dir's name. Use ex: \"test dir\"") //отправляем ошибку клиенту
		return
	}
	err := os.Mkdir(rootPath + path,0777) //0777 - максимальный уровень доступа к папке (полные права на чтение\запись). Можно регулировать http://www.rhd.ru/docs/manuals/enterprise/RHEL-AS-2.1-Manual/getting-started-guide/s1-navigating-chmodnum.html
	if err != nil {
		io.WriteString(s, err.Error()) //отправляем ошибку клиенту
	}
}