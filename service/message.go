package service

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/lw396/WeComCopilot/internal/errors"
	mysql "github.com/lw396/WeComCopilot/internal/repository/gorm"
	"github.com/lw396/WeComCopilot/internal/repository/sqlite"
	"github.com/lw396/WeComCopilot/pkg/db"
	"github.com/lw396/WeComCopilot/pkg/util"
	"gorm.io/gorm"
)

type MessageInfo struct {
	UserName string `json:"user_name"`
	Seq      uint64 `json:"seq"`
	DBName   string `json:"db_name"`
}

func (a *Service) ScanMessage(ctx context.Context, userName string) (result *MessageInfo, err error) {
	var dbName string
	var seq *sqlite.SQLiteSequence
	var name string = "Chat_" + hex.EncodeToString(util.Md5([]byte(userName)))
	for i := 0; i < 10; i++ {
		dbName = fmt.Sprintf(sqlite.MessageDB, i)
		var tx *gorm.DB
		if tx, err = a.sqlite.OpenDB(ctx, dbName); err != nil {
			return
		}
		if seq, err = a.sqlite.CheckMessageExistDB(ctx, tx, name); err != nil {
			if !db.IsRecordNotFound(err) {
				return
			}
			continue
		}
		break
	}
	if seq == nil {
		err = errors.New(errors.CodeDB, "not found message")
		return
	}
	result = &MessageInfo{
		DBName:   dbName,
		UserName: userName,
		Seq:      seq.Seq,
	}
	return
}

func (a *Service) SaveMessageContent(ctx context.Context, data *GroupContact) (err error) {
	tx, err := a.sqlite.OpenDB(ctx, data.DBName)
	if err != nil {
		return
	}
	msgName := "Chat_" + hex.EncodeToString(util.Md5([]byte(data.UsrName)))
	err = a.sqlite.BindMessage(ctx, tx, data.DBName, msgName)
	if err != nil {
		return
	}
	message, err := a.sqlite.GetMessageContent(ctx, data.DBName, msgName)
	if err != nil {
		return
	}

	if err = a.rep.SaveGroupContact(ctx, &mysql.GroupContact{
		UsrName:         data.UsrName,
		Nickname:        data.Nickname,
		HeadImgUrl:      data.HeadImgUrl,
		ChatRoomMemList: data.ChatRoomMemList,
		DBName:          data.DBName,
	}); err != nil {
		return
	}

	fmt.Println(message[0])
	return
}
