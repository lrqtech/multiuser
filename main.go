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
===========
 Menu List
     _
   V1.1
===========
1. Find a user (By id)
2. Find a user (By username)
3. Add a user
4. Chang password
5. Delete a user
6. Exit`

type Logs struct {
	Id       int    `xorm:"pk autoincr INTEGER"`
	Type     string `xorm:"TEXT"`
	Log      string `xorm:"TEXT"`
	Addtime  string `xorm:"TEXT"`
	Uid      int    `xorm:"default '1' integer"`
	Username string `xorm:"default 'system' TEXT"`
}

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

func AddLog(log string) error {
	lt := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
	_, err := engine.InsertOne(&Logs{Type: "用户管理", Log: log, Addtime: lt, Username: "system"})
	return err
}

func main() {
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
	log := "有人启动了[面板多用户管理（非官方）]"
	_ = AddLog(log)
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
			log := fmt.Sprintf("[面板多用户管理（非官方）]-有人增加用户[ %s ]", name)
			_ = AddLog(log)
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
			log := fmt.Sprintf("[面板多用户管理（非官方）]-有人修改了用户(ID)[ %d ]的密码", id)
			_ = AddLog(log)
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
				log := fmt.Sprintf("[面板多用户管理（非官方）]-有人删除了用户(ID)[ %d ]", id)
				_ = AddLog(log)
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
