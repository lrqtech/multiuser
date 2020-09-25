# 宝塔面板多用户管理-三方插件

- 功能介绍：支持管理面板多用户，可查看各个用户日志，但不支持权限分配
- 支持版本：6.8.4及以上（Linux Only）
- 安装方法：SSH输入下面指令


（golang1.15的安装和go mod的配置自行百度）



编译：

    git clone https://github.com/lrqtech/multiuser.git
    go mod tidy
    go build


运行：

    ./cli

