package gameConfig

import (
	"encoding/json"
	"log"
	"os"
)

const (
	TypeDurable        = 1 //耐久
	TypeCost           = 2
	TypeArmyTeams      = 3 //队伍数量
	TypeSpeed          = 4 //速度
	TypeDefense        = 5 //防御
	TypeStrategy       = 6 //谋略
	TypeForce          = 7 //攻击武力
	TypeConscriptTime  = 8 //征兵时间
	TypeReserveLimit   = 9 //预备役上限
	TypeUnkonw         = 10
	TypeHanAddition    = 11
	TypeQunAddition    = 12
	TypeWeiAddition    = 13
	TypeShuAddition    = 14
	TypeWuAddition     = 15
	TypeDealTaxRate    = 16 //交易税率
	TypeWood           = 17
	TypeIron           = 18
	TypeGrain          = 19
	TypeStone          = 20
	TypeTax            = 21 //税收
	TypeExtendTimes    = 22 //扩建次数
	TypeWarehouseLimit = 23 //仓库容量
	TypeSoldierLimit   = 24 //带兵数量
	TypeVanguardLimit  = 25 //前锋数量
)

const (
	Main          = 0  //主城
	JiaoChang     = 13 //校场
	TongShuaiTing = 14 //统帅厅
	JiShi         = 15 //集市
	MBS           = 16 //募兵所
)

type conditions struct {
	Type  int `json:"type"`
	Level int `json:"level"`
}

type facility struct {
	Title      string       `json:"title"`
	Des        string       `json:"des"`
	Name       string       `json:"name"`
	Type       int8         `json:"type"`
	Additions  []int8       `json:"additions"`
	Conditions []conditions `json:"conditions"`
	Levels     []fLevel     `json:"levels"`
}

type NeedRes struct {
	Decree int `json:"decree"`
	Grain  int `json:"grain"`
	Wood   int `json:"wood"`
	Iron   int `json:"iron"`
	Stone  int `json:"stone"`
	Gold   int `json:"gold"`
}

type fLevel struct {
	Level  int     `json:"level"`
	Values []int   `json:"values"`
	Need   NeedRes `json:"need"`
	Time   int     `json:"time"` //升级需要的时间
}

type conf struct {
	Name string
	Type int8
}

type facilityConf struct {
	Title     string `json:"title"`
	List      []conf `json:"list"`
	facilitys map[int8]*facility
}

var FacilityConf = &facilityConf{}

const facilityFile = "/conf/game/facility/facility.json"
const facilityPath = "/conf/game/facility/"

func (f *facilityConf) Load() {
	//获取当前文件路径
	currentDir, _ := os.Getwd()
	//配置文件位置
	cf := currentDir + facilityFile
	cfPath := currentDir + facilityPath
	//打包后 程序参数加入配置文件路径
	if len(os.Args) > 1 {
		if path := os.Args[1]; path != "" {
			cf = path + facilityFile
			cfPath = path + facilityPath
		}
	}
	data, err := os.ReadFile(cf)
	if err != nil {
		log.Println("城池设施读取失败")
		panic(err)
	}
	err = json.Unmarshal(data, f)
	if err != nil {
		log.Println("城池设施格式定义失败")
		panic(err)
	}
	f.facilitys = make(map[int8]*facility, len(f.List))

	files, err := os.ReadDir(cfPath)
	if err != nil {
		log.Println("设施文件读取失败")
		panic(err)
	}
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		name := v.Name()
		if name == "facility.json" {
			continue
		}
		data, err := os.ReadFile(cfPath + v.Name())
		if err != nil {
			panic(err)
		}
		fac := &facility{}
		err = json.Unmarshal(data, fac)
		if err != nil {
			panic(err)
		}
		f.facilitys[fac.Type] = fac
	}
}

func (f *facilityConf) CostTime(fType int8, level int8) int {
	if level <= 0 {
		return 0
	}
	fa, ok := f.facilitys[fType]
	if ok {
		if int8(len(fa.Levels)) >= level {
			return fa.Levels[level-1].Time - 2 //比客户端快2s，保证客户端倒计时完一定是升级成功了
		} else {
			return 0
		}
	} else {
		return 0
	}
}

func (f *facilityConf) GetValues(fType int8, level int8) []int {
	if level <= 0 {
		return []int{}
	}

	fa, ok := f.facilitys[fType]
	if ok {
		if int8(len(fa.Levels)) >= level {
			return fa.Levels[level-1].Values
		} else {
			return []int{}
		}
	} else {
		return []int{}
	}
}

func (f *facilityConf) GetAdditions(fType int8) []int8 {
	fa, ok := f.facilitys[fType]
	if ok {
		return fa.Additions
	} else {
		return []int8{}
	}
}

func (f *facilityConf) MaxLevel(fType int8) int {
	fa, ok := f.facilitys[fType]
	if ok {
		return len(fa.Levels)
	} else {
		return 0
	}
}

func (f *facilityConf) Need(fType int8, level int8) NeedRes {
	fa, ok := f.facilitys[fType]
	if ok {
		return fa.Levels[level].Need
	}
	return NeedRes{}
}
