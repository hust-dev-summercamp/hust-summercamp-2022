package db

import (
	"fmt"
	mydb "summercamp-filestore/db/mysql"
	"time"
)

// UserFile : 用户文件表结构体
type UserFile struct {
	UserName    string `json:"user_name"`
	FileHash    string `json:"file_hash"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	UploadAt    string `json:"upload_at"`
	LastUpdated string `json:"last_updated"`
}

// OnUserFileUploadFinished : 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	fmt.Println("inserting user file")
	stmt, err := mydb.DBConn().Prepare(
		"insert into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`,`status`) values (?,?,?,?,?,1)")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		return false
	}
	return true
}

// QueryUserFileMetas : 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at," +
			"last_update from tbl_user_file where user_name=? and status!=2 limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	fmt.Println(stmt)
	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize,
			&ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}

// RenameFileName : 文件重命名
func RenameFileName(username, filehash, filename string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user_file set file_name=? where user_name=? and file_sha1=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// DeleteUserFile : 删除文件(标记删除)
func DeleteUserFile(username, filehash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user_file set status=2 where user_name=? and file_sha1=? and status =1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// QueryUserFileMeta : 获取用户单个文件信息
func QueryUserFileMeta(username string, filehash string) (*UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at," +
			"last_update from tbl_user_file where user_name=? and file_sha1=? and status=1 limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		return nil, err
	}

	ufile := UserFile{}
	if rows.Next() {
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize,
			&ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}

	return &ufile, nil
}

// IsUserFileUploaded : 用户文件是否已经上传过
func IsUserFileUploaded(username string, filehash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select 1 from tbl_user_file where user_name=? and file_sha1=? and status=1 limit 1")
	rows, err := stmt.Query(username, filehash)
	if err != nil {
		return false
	} else if rows == nil || !rows.Next() {
		return false
	}
	return true
}
