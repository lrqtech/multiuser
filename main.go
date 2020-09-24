package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os/exec"
	"time"
	"xorm.io/xorm"
)

const prompt = `
Menu List
1. Find a user (By id)
2. Find a user (By username)
3. Add a user
4. Chang password
5. Delete a user
6. Exit`

type Users struct {
	Id        int    `xorm:"pk autoincr INTEGER"`
	Username  string `xorm:"TEXT"`
	Password  string `xorm:"TEXT"`
	LoginIp   string `xorm:"TEXT"`
	LoginTime string `xorm:"TEXT"`
	Phone     string `xorm:"TEXT"`
	Email     string `xorm:"TEXT"`
}

var engine *xorm.Engine

func FindUserId(id int) (*Users, error) {
	results := &Users{}
	has, err := engine.ID(id).Get(results)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("no data")
	}
	return results, nil
}

func FindUserName(name string) (*Users, error) {
	results := &Users{Username: name}
	has, err := engine.Get(results)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("no data")
	}
	return results, nil
}

func AddUser(username, passwd string) error {
	lt := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
	_, err := engine.InsertOne(&Users{Username: username, Password: passwd, LoginTime: lt})
	return err
}

func ChangPasswd(id int, passwd string) error {
	results := &Users{Password: passwd}
	_, err := engine.ID(id).Update(results)
	return err
}

func DelUser(id int) error {
	results := &Users{}
	_, err := engine.ID(id).Delete(results)
	return err
}

func Md5(str string) string {
	data := []byte(str)
	e := md5.Sum(data)
	result := fmt.Sprintf("%x", e)
	return result
}

func main() {
	fmt.Println("Loading......")
	cmds := exec.Command("/bin/bash", "-c", `chmod 777 /www/server/panel/data/default.db`)
	if err := cmds.Start(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	var err error
	engine, err = xorm.NewEngine("sqlite3", "/www/server/panel/data/default.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := engine.Ping(); err != nil {
		fmt.Println("Ping database fail")
		return
	}
	fmt.Println("Welcome to CLI")
Exit:
	for {
		fmt.Println(prompt)
		var num int
		fmt.Scanf("%d \n", &num)
		switch num {
		case 1:
			fmt.Println("Please enter user id to find")
			var id int
			fmt.Scanf("%d \n", &id)
			results, err := FindUserId(id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Username: %s \n", results.Username)
				fmt.Printf("Create time: %s \n", results.LoginTime)
			}
		case 2:
			fmt.Println("Please enter user name to find")
			var name string
			fmt.Scanf("%s \n", &name)
			results, err := FindUserName(name)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Username: %d \n", results.Id)
				fmt.Println("Now you can find this user by option 1 ")
			}
		case 3:
			fmt.Println("Adding user......")
			fmt.Println("Please enter new user name")
			var name string
			fmt.Scanf("%s \n", &name)
			fmt.Println("Please enter password for new user")
			var passwd string
			fmt.Scanf("%s \n", &passwd)
			err = AddUser(name, Md5(passwd))
			if err != nil {
				fmt.Println(err)
			}
		case 4:
			fmt.Println("Please enter user ID")
			var id int
			fmt.Scanf("%d \n", &id)
			fmt.Println("Please enter new password")
			var passwd string
			fmt.Scanf("%s \n", &passwd)
			err = ChangPasswd(id, passwd)
			if err != nil {
				fmt.Println(err)
			}
		case 5:
			fmt.Println("Please enter user id that you want to delete")
			var id int
			fmt.Scanf("%d \n", &id)
			if id == 1 {
				fmt.Println("You can not delete the default user")
			} else {
				err = DelUser(id)
				if err != nil {
					fmt.Println(err)
				}
			}
		case 6:
			break Exit
		}
	}
	cmde := exec.Command("/bin/bash", "-c", `chmod 600 /www/server/panel/data/default.db`)
	_ = cmde.Start()
	fmt.Println("Goodbye")
}
