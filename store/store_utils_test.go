package store

import (
	"strings"
	"testing"
)

func Test_FileName2Offset(t *testing.T) {
	fileName := "00201900000000000200"
	offset, err := FileName2Offset(fileName)
	if err != nil {
		t.Error(err)
	}
	if strings.EqualFold(fileName, Offset2FileName(offset)) == false {
		t.Errorf("filename %s 与 offset %s 相互间转换异常", fileName, Offset2FileName(offset))
	}
}
