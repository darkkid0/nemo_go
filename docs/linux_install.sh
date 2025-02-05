#!/bin/bash

#Tested on ubuntu18.04 LTS

#apt-get
#只有server需要安装 mysql-server和rabbitmq-server
#只有worker需要安装 nmap masscan chromium-browser
    && sudo apt-get install wget curl vim net-tools git unzip \
    mysql-server rabbitmq-server \
    nmap masscan chromium-browser --fix-missing

# init databse and rabbitmq
# 只有server需要
# 由于mysql5.7版本后默认安装后不能使用root和空密码在本机登录，
# 系统自动生成的用户名密码位于/etc/mysql/debian.cnf，请替换并使用配置文件中的user和password：
sudo service mysql start \
    && mysql -u root -e 'CREATE DATABASE `nemo` DEFAULT CHARACTER SET utf8mb4;' \
    && mysql -u root -e 'CREATE USER "nemo"@"%" IDENTIFIED BY "nemo2020";GRANT ALL PRIVILEGES ON nemo.* TO "nemo"@"%";FLUSH PRIVILEGES;' \
    && mysql -u root nemo < ../docker/mysql/initdb.d/docs/nemo.sql \

