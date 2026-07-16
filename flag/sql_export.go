package flag

import (
	"fmt"
	"os"
	"os/exec"
	"server/global"
	"time"
)

// SQLExport 导出 MySQL 数据
func SQLExport() error {
	mysql := global.Configs.Mysql

	timer := time.Now().Format("20060102")
	sqlPath := fmt.Sprintf("mysql_%s.sql", timer)
	//准备一个外部命令，docker exec mysql mysqldump -u<用户名> -p<密码> <数据库名>进行表迁移
	cmd := exec.Command("docker", "exec", "mysql", "mysqldump", "-u"+mysql.Username, "-p"+mysql.Password, mysql.DBName)

	outFile, err := os.Create(sqlPath)
	if err != nil {
		return err
	}
	defer outFile.Close() //关闭文件

	cmd.Stdout = outFile //Stdout 是io.Writer类型，让它等于刚刚创建的文件即可，确保了数据库备份的内容被直接保存为本地文件，而不是在控制台刷屏。
	return cmd.Run()
}
