package services

import (
	"github.com/sirupsen/logrus"
	"sheetServerApi/global"
	"sheetServerApi/internal/middlewares/constants"
	model "sheetServerApi/internal/model/db"
	"sheetServerApi/internal/model/params"
	"sheetServerApi/internal/utils"
	"strconv"
	"time"
)

/**
* @program: src
* @description: v1版本的报表生成
* @author: 占翔昊
* @create 2020-10-14 14:26
**/

// 最新文件的生成
var sheetname = "Sheet1"
// sqlx连接器
var opsqlx model.OpSqlxExcelMetaDao
// gorm连接器
var opgorm model.OpGormExcelMetaDao

// 生成文件的主逻辑
func GenerateSheetFile(req params.SheetParamsReq) (string,error) {
	/**
	创建文件
	 */
	f,err := utils.Init_file(sheetname)
	if err != nil {
		logrus.Fatal(err)
		return "error",err
	}

	/**
	设置行列宽度
	*/
	if err := SetColsAndRowslength(f,req.Cell,sheetname);err!=nil{
		logrus.Fatal(err)
		return "error",err
	}

	/**
	循环写入样式元数据和单元格值
	 */
	if err := SetBlockStyleAndValue(f,req.Data,sheetname,req.Author);err!=nil {
		logrus.Fatal(err)
		return "error",err
	}

	// 协程写入数据库，这里要添加异常处理
	var OpExcelDao model.OpGormExcelMetaDao
	go OpExcelDao.WriteData(global.DBOrmEngine,req)

	file_name := req.Name + constants.Name_time_mark + strconv.Itoa(int(time.Now().Unix())) + ".xlsx"
	file_dir := global.AppSetting.ExcelFileDir + file_name
	// 文件写入磁盘
	//file_dir := req.Name + constants.Name_time_mark +  strconv.Itoa(int(time.Now().Unix())) + ".xlsx"
	if err := f.SaveAs(file_dir);err!=nil {
		logrus.Fatal(err)
		return "error",err
	}
	url := global.AppSetting.ReleaseUrl + file_name
	return url ,nil
}



// 获取数据库表格信息
func GetTableMetaInfo(tableName string) ([]params.SheetTableMetaInfo,error) {
	// 首先去判断数据库字段是否正确
	data,err := opsqlx.GetTableMetaInfo(global.DBSqlxEngine,tableName)
	if err!=nil {
		logrus.Fatal(err)
		return nil,err
	}
	return data,nil
}

// 获取所有原始字段
func GetExcelRawDatas(id int64) (string,error)  {
	data,err := opsqlx.GetSheetRawData(global.DBSqlxEngine,id)
	if err != nil {
		logrus.Fatal(err)
		return data , err
	}
	return data,nil
}

func GetSheetHistory(req params.SheetHistoryReq) ([]params.SheetHistoryInfo,error){
	data,err := opsqlx.GetSheetHistory(global.DBSqlxEngine,req)
	if err != nil {
		logrus.Fatal(err)
		return data , err
	}
	return data,nil
}