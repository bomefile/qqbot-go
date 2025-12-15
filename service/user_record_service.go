package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/dao"
)

type UserRecordInput struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Concat string `json:"concat"`
}

type UserRecordUpdateInput struct {
	UID    string  `json:"uid"`
	Name   *string `json:"name,omitempty"`
	Concat *string `json:"concat,omitempty"`
}

func UserRecordHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	if r.Method == http.MethodGet {
		uid := r.URL.Query().Get("uid")
		pageNumStr := r.URL.Query().Get("page_num")
		pageSizeStr := r.URL.Query().Get("page_size")
		pageNum := 1
		pageSize := 10
		if v, err := strconv.Atoi(pageNumStr); err == nil && v > 0 {
			pageNum = v
		}
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 {
			pageSize = v
		}
		records, err := dao.ListUserRecords(uid, pageNum, pageSize)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			totalSize, err := dao.CountUserRecords(uid)
			if err != nil {
				res.Code = -1
				res.ErrorMsg = err.Error()
			} else {
				totalPage := 0
				if pageSize > 0 {
					totalPage = int((totalSize + int64(pageSize) - 1) / int64(pageSize))
				}
				res.Data = map[string]interface{}{
					"page_num":   pageNum,
					"page_size":  pageSize,
					"total_page": totalPage,
					"total_size": totalSize,
					"list":       records,
				}
			}
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

func UserRecordAddHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	if r.Method != http.MethodPost {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	} else {
		var in UserRecordInput
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&in); err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else if in.UID == "" || in.Name == "" {
			res.Code = -1
			res.ErrorMsg = "缺少 uid 或 name"
		} else if err := dao.CreateUserRecordIfNotExists(in.UID, in.Name, in.Concat); err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = "ok"
		}
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

func UserRecordUpdateHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	if r.Method != http.MethodPost {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	} else {
		var in UserRecordUpdateInput
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&in); err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else if in.UID == "" {
			res.Code = -1
			res.ErrorMsg = "缺少 uid"
		} else if err := dao.UpdateUserRecordByUID(in.UID, in.Name, in.Concat); err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = "ok"
		}
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

type UserRecordCheckinInput struct {
	UID string `json:"uid"`
}

func UserRecordCheckinHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	if r.Method != http.MethodPost {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	} else {
		var in UserRecordCheckinInput
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&in); err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else if in.UID == "" {
			res.Code = -1
			res.ErrorMsg = "缺少 uid"
		} else if err := dao.CheckinUserRecord(in.UID); err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = "ok"
		}
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}
