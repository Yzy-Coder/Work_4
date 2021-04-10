package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//保存连接数据库的常量
const (
	userName     = "root"
	passWord     = ""
	host         = "127.0.0.1"
	port         = "3306"
	dataBaseName = "godatabase"
)

//初始化数据库，返回 *sql.DB类型的数据库操作句柄
func DataBaseInit() *sql.DB {

	//获得数据库操作的句柄
	DB, _ := sql.Open("mysql", userName+":"+passWord+"@("+host+":"+port+")/"+dataBaseName)

	//测试连接
	err := DB.Ping()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("连接成功")
	}

	//返回句柄
	return DB
}

//用户注册操作，传入用户名，用户密码，如果成功返回 用户ID UserId，如果错误返回 -1
func Registration(DB *sql.DB, UserName string, UserPassWord string) int64 {

	//判断密码是否为空，如果为空直接返回-1
	if UserPassWord == "" {

		fmt.Println("密码不能为0")
		return -1

	}

	//准备检查用户名是否存在操作的SQL语句
	querySql := " SELECT COUNT(UserName) AS COUNT FROM UserInfo WHERE UserName = ? "

	//执行querySql语句，获得返回的结果
	queryQuery, err := DB.Query(querySql, UserName)

	//如果 err!=nil 返回-1
	if err != nil {

		fmt.Println(err)
		return -1
	}

	//判断是否存在下一行（因为第一次执行，即获得第一行数据）
	queryQuery.Next()
	var COUNT int64
	queryQuery.Scan(&COUNT)

	//如果 COUNT > 0 ，则说明已存在同样的用户名，返回 -1
	if COUNT > 0 {

		fmt.Println("用户名重复")
		return -1
	}

	//准备插入用户信息的SQL语句
	insertSql := " INSERT INTO UserInfo ( UserName , UserPassWord )  VALUES ( ? , ? )"

	//获得插入的结果，并通过insertResult来获取
	insertResult, err := DB.Exec(insertSql, UserName, UserPassWord)

	//如果 err!=nil 返回-1
	if err != nil {

		fmt.Println(err)
		return -1
	}

	//获得被插入用户操作所影响的行数， 如果rowsAffected == 0 ，则说明插入失败
	rowsAffected, err := insertResult.RowsAffected()
	if err != nil || rowsAffected == 0 {

		fmt.Println(err)
		fmt.Println("注册失败")
		return -1
	}

	//利用Result.LastInsertId()获得 UserId
	lastInserId, err := insertResult.LastInsertId()
	if err != nil {

		fmt.Println(err)
		return -1
	}

	//返回UserId
	return lastInserId

}

//用户登录操作，传入用户名，用户密码，如果成功返回 1 ，如果错误返回-1，如果用户名或密码错误返回 -2
func SignIn(DB *sql.DB, UserName string, UserPassWord string) int64 {

	//准备查询该用户名的密码的SQL语句
	querySql := "SELECT UserPassWord FROM UserInfo WHERE UserName = ?"

	//执行查询语句获得返回的行
	queryRows, err := DB.Query(querySql, UserName)
	if err != nil {
		fmt.Println("数据库查询出错", err)
		return -1
	}

	//定义一个userPassWord来存放取出来的密码
	var userPassWord string

	//如果取出来了数据，则将其赋值给userPassWord
	if queryRows.Next() {
		queryRows.Scan(&userPassWord)
	}

	//如果账号密码正确，则返回 1
	if UserPassWord == userPassWord {
		fmt.Println("登录成功")
		return 1
	}

	//如果密码错误，说明账号或密码错误 返回 -2
	if UserPassWord != userPassWord {
		fmt.Println("账号或密码错误")
		return -2
	}

	return -1
}

//博客发布操作，传入用户Id UserId，博客标题 BlogName，博客内容 BlogContent，如果成功返回博客Id BlogId，如果失败则返回-1
func Publish(DB *sql.DB, UserId string, BlogName string, BlogContent string) int64 {

	//准备查询是否存在该用户存在的SQL语句
	querySql := "SELECT COUNT(UserId) FROM UserInfo WHERE UserId = ?"

	//获得被查询语句返回的行并赋值给count，即是否存在该用户，如果count == 0 则说明该用户不存在
	var count int
	QueryRows, err := DB.Query(querySql, UserId)
	if err != nil {

		fmt.Println(err)
		return -1
	}
	QueryRows.Next()
	QueryRows.Scan(&count)
	if count == 0 {
		fmt.Println(err, "用户不存在")
		return -1
	}

	//准备执行博客插入的SQL语句
	insertSql := "INSERT INTO BlogInfo ( UserId , BlogName , BlogContent ) VALUES ( ? , ? , ? )"

	//获得插入的结果，将被插入操作影响的行数赋值给rowsAffected
	insertResult, _ := DB.Exec(insertSql, UserId, BlogName, BlogContent)
	rowsAffected, err := insertResult.RowsAffected()

	//如果rowsAffected <= 0 ， 则插入失败 ，返回-1
	if err != nil || rowsAffected <= 0 {
		fmt.Println(err, "插入失败")
		return -1
	}

	//获得该博客的Id BlogId并返回
	BlogId, err := insertResult.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return BlogId

	////准备查询博客Id的SQL语句
	//querySql = "SELECT BlogId FROM BlogInfo WHERE UserId = ? ORDER BY BlogId DESC "
	//
	////执行查询语句,获得BlogId
	//var BlogId int
	//QueryRows , err = DB.Query(querySql,UserId)
	//if err != nil {
	//	fmt.Println(err)
	//	return -1
	//}
	//QueryRows.Next()
	//QueryRows.Scan(&BlogId)
	//
	//fmt.Println(BlogId)
	////返回博客ID BlogId
	//return int64(BlogId)
}

//博客点赞操作，传入用户Id UserId，博客Id BlogId ， 如果成功则返回该博客的点赞数 BlogThumbsUpNum ，失败则返回 -1
func ThumbsUp(DB *sql.DB, UserId string, BlogId string) int64 {

	//准备查询是否存在该用户存在的SQL语句
	querySql := "SELECT COUNT(UserId) FROM UserInfo WHERE UserId = ?"
	//执行查询语句并获得查询语句返回的行并赋值给count，即是否存在该用户，如果count == 0 则说明该用户不存在
	var count int
	QueryRows, err := DB.Query(querySql, UserId)
	if err != nil {

		fmt.Println(err)
		return -1
	}
	QueryRows.Next()
	QueryRows.Scan(&count)
	if count == 0 {
		fmt.Println(err, "用户不存在")
		return -1
	}

	//准备查询博客是否存在的SQL语句
	querySql = "SELECT COUNT(BlogId) FROM BlogInfo WHERE BlogId = ? "
	//执行查询语句并获得查询语句返回的行并赋值给count, 即是否存在该博客 , 如果count == 0 则说明该博客不存在
	QueryRows, err = DB.Query(querySql, BlogId)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	QueryRows.Next()
	QueryRows.Scan(&count)
	if count == 0 {
		fmt.Println("博客不存在")
		return -1
	}

	//准备查询是否存在该用户与该博客的点赞依赖SQL语句
	querySql = "SELECT COUNT(*) FROM BlogThumbsUpInfo WHERE UserId=? AND BlogId=?"
	//获得被查询语句返回的行并赋值给count，即是否存在该依赖，如果count != 0 则说明该依赖已存在，不能再点赞
	QueryRows, err = DB.Query(querySql, UserId, BlogId)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	QueryRows.Next()
	QueryRows.Scan(&count)
	if count != 0 {
		fmt.Println("无法重复点赞")
		return -1
	}

	//准备插入点赞依赖关系的SQL语句
	insertSql := "INSERT INTO BlogThumbsUpInfo (UserId , BlogId) VALUES (?,?)"
	//执行插入语句，并将受影响的行赋值给rowsAffected ， 如果 err != nil 或 rowsAffected == 0 则说明插入失败， 返回-1
	insertResult, err := DB.Exec(insertSql, UserId, BlogId)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	rowsAffected, err := insertResult.RowsAffected()
	if err != nil || rowsAffected == 0 {
		fmt.Println(err, "rowsAffected=", rowsAffected)
		return -1
	}

	//准备执行更新该博客的点赞数的SQL语句
	updateSql := "UPDATE BlogInfo SET BlogThumbsUpNum = BlogThumbsUpNum + 1 WHERE BlogId = ?"
	//执行更新语句并将受影响的行赋值给rowsAffected，如果 err != nil 或 rowsAffected == 0 ，说明更新失败 ，返回 -1
	updataResult, err := DB.Exec(updateSql, BlogId)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	rowsAffected, err = updataResult.RowsAffected()
	if err != nil || rowsAffected == 0 {
		fmt.Println(err, "rowsAffected=", rowsAffected)
		return -1
	}

	//准备查询该博客现在点赞数的SQL语句
	querySql = "SELECT BlogThumbsUpNum FROM BlogInfo WHERE BlogId = ?"
	//执行查询语句并将返回的Rows赋值给queryRows
	queryRows, err := DB.Query(querySql, BlogId)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	//存在查询返回的行，则返回该博客现在的点赞数ThisThumbsUpNum
	if queryRows.Next() {

		var ThisThumbsUpNum int64
		queryRows.Scan(&ThisThumbsUpNum)
		return ThisThumbsUpNum
	}

	return -1
}

//func main() {
//
//	DB := DataBaseInit()
//	//Registration(DB,"chy","516")
//	//SignIn(DB,"Yzy","123456")
//	//fmt.Println(Publish(DB, 1, "Test", "This is a test blog"))
//	//fmt.Println(ThumbsUp(DB,1,44))
//
//}
