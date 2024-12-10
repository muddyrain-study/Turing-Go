package data

import (
	"Turing-Go/server/game/model"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

const (
	UnionDismiss = 0 //解散
	UnionRunning = 1 //运行中
)

type Coalition struct {
	Id           int       `xorm:"id pk autoincr"`
	Name         string    `xorm:"name"`
	Members      string    `xorm:"members"`
	MemberArray  []int     `xorm:"-"`
	CreateId     int       `xorm:"create_id"`
	Chairman     int       `xorm:"chairman"`
	ViceChairman int       `xorm:"vice_chairman"`
	Notice       string    `xorm:"notice"`
	State        int8      `xorm:"state"`
	Ctime        time.Time `xorm:"ctime"`
}

func (c *Coalition) TableName() string {
	return "coalition"
}
func (c *Coalition) AfterSet(name string, cell xorm.Cell) {
	if name == "members" {
		if cell != nil {
			ss, ok := (*cell).([]uint8)
			if ok {
				err := json.Unmarshal(ss, &c.MemberArray)
				if err != nil {
					log.Println("coalition AfterSet error", err)
				}
			}
			if c.MemberArray == nil {
				c.MemberArray = make([]int, 0)
			}
		}
	}
}

func (c *Coalition) ToModal() interface{} {
	u := model.Union{}
	u.Name = c.Name
	u.Id = c.Id
	u.Notice = c.Notice
	u.Cnt = c.Cnt()
	return u
}

func (c *Coalition) Cnt() int {
	return len(c.MemberArray)
}

type CoalitionApply struct {
	Id      int       `xorm:"id pk autoincr"`
	UnionId int       `xorm:"union_id"`
	RId     int       `xorm:"rid"`
	State   int8      `xorm:"state"`
	Ctime   time.Time `xorm:"ctime"`
}

func (c *CoalitionApply) TableName() string {
	return "coalition_apply"
}

const (
	UnionOpCreate    = 0 //创建
	UnionOpDismiss   = 1 //解散
	UnionOpJoin      = 2 //加入
	UnionOpExit      = 3 //退出
	UnionOpKick      = 4 //踢出
	UnionOpAppoint   = 5 //任命
	UnionOpAbdicate  = 6 //禅让
	UnionOpModNotice = 7 //修改公告
)

type CoalitionLog struct {
	Id       int       `xorm:"id pk autoincr"`
	UnionId  int       `xorm:"union_id"`
	OPRId    int       `xorm:"op_rid"`
	TargetId int       `xorm:"target_id"`
	State    int8      `xorm:"state"`
	Des      string    `xorm:"des"`
	Ctime    time.Time `xorm:"ctime"`
}

func (c *CoalitionLog) TableName() string {
	return "coalition_log"
}
