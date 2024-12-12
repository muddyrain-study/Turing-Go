package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/global"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"Turing-Go/utils"
	"github.com/go-xorm/xorm"
	"log"
	"math/rand"
	"sync"
	"time"
)

var RoleCityService = &roleCityService{
	posCity:  make(map[int]*data.MapRoleCity),
	roleCity: make(map[int][]*data.MapRoleCity),
	mutex:    sync.RWMutex{},
	dbCity:   make(map[int]*data.MapRoleCity),
}

type roleCityService struct {
	posCity  map[int]*data.MapRoleCity
	roleCity map[int][]*data.MapRoleCity
	mutex    sync.RWMutex
	dbCity   map[int]*data.MapRoleCity
}

func (r *roleCityService) Load() {
	err := db.Engine.Find(&r.dbCity)
	dbCity := r.dbCity
	if err != nil {
		log.Println("RoleCityService load role_city table error")
		return
	}
	//转成posCity、roleCity
	for _, v := range dbCity {
		posId := global.ToPosition(v.X, v.Y)
		r.posCity[posId] = v
		_, ok := r.roleCity[v.RId]
		if ok == false {
			r.roleCity[v.RId] = make([]*data.MapRoleCity, 0)
		}
		r.roleCity[v.RId] = append(r.roleCity[v.RId], v)
	}
	//耐久度计算 后续做
}

func (r *roleCityService) InitCity(rid int, nickName string, req *net.WsMsgReq) error {
	session := req.Context.Get("session").(*xorm.Session)
	roleCity := data.MapRoleCity{}
	ok, err := db.Engine.Table(roleCity).Where("rid=?", rid).Get(&roleCity)
	if err != nil {
		log.Println("查询角色城池异常", err)
		return common.New(constant.DBError, "数据库错误")
	}
	if !ok {
		for {
			roleCity.X = rand.Intn(global.MapWidth)
			roleCity.Y = rand.Intn(global.MapHeight)
			if r.isCanBuild(roleCity.X, roleCity.Y) {
				roleCity.RId = rid
				roleCity.Name = nickName
				roleCity.CurDurable = gameConfig.Basic.City.Durable
				roleCity.CreatedAt = time.Now()
				roleCity.IsMain = 1
				if session != nil {
					_, err = session.Table(roleCity).Insert(&roleCity)
				} else {
					_, err = db.Engine.Table(roleCity).Insert(&roleCity)
				}
				if err != nil {
					log.Println("插入角色城池异常", err)
					return common.New(constant.DBError, "数据库错误")
				}
				err := session.Commit()
				if err != nil {
					log.Println("插入角色城池异常", err)
					return common.New(constant.DBError, "数据库错误")
				}
				posId := global.ToPosition(roleCity.X, roleCity.Y)
				r.posCity[posId] = &roleCity
				_, ok := r.roleCity[rid]
				if !ok {
					r.roleCity[rid] = make([]*data.MapRoleCity, 0)
				} else {
					r.roleCity[rid] = append(r.roleCity[rid], &roleCity)
				}
				r.dbCity[roleCity.CityId] = &roleCity
				err = CityFacilityService.TryCreate(roleCity.CityId, rid, session)
				if err != nil {
					log.Println("插入城市设施异常", err)
					return common.New(constant.DBError, "数据库错误")
				}
				break
			}
		}

	}
	return nil

}

func (r *roleCityService) isCanBuild(x int, y int) bool {
	confs := gameConfig.MapRes.Confs
	pIndex := global.ToPosition(x, y)
	_, ok := confs[pIndex]
	if !ok {
		return false
	}
	sysBuild := gameConfig.MapRes.SysBuild
	//系统城池的5格内 不能创建玩家城池
	for _, v := range sysBuild {
		if v.Type == gameConfig.MapBuildSysCity {
			if x >= v.X-5 &&
				x <= v.X+5 &&
				y >= v.Y-5 &&
				y <= v.Y+5 {
				return false
			}
		}
	}
	//玩家城池的5格内 也不能创建城池
	for i := x - 5; i <= x+5; i++ {
		for j := y - 5; j <= y+5; j++ {
			posId := global.ToPosition(i, j)
			_, ok := r.posCity[posId]
			if ok {
				return false
			}
		}
	}
	return true
}

func (r *roleCityService) GetRoleCities(rid int) ([]model.MapRoleCity, error) {
	cities := make([]data.MapRoleCity, 0)
	city := &data.MapRoleCity{}
	err := db.Engine.Table(city).Where("rid=?", rid).Find(&cities)
	modelCities := make([]model.MapRoleCity, len(cities))

	if err != nil {
		log.Println("查询角色城池异常", err)
		return modelCities, common.New(constant.DBError, "数据库错误")
	}
	for i, roleCity := range cities {
		modelCities[i] = roleCity.ToModel().(model.MapRoleCity)
	}
	return modelCities, nil
}

func (r *roleCityService) ScanBlock(req *model.ScanBlockReq) ([]model.MapRoleCity, error) {
	x := req.X
	y := req.Y
	length := req.Length
	if x < 0 || x >= global.MapWidth || y < 0 || y >= global.MapHeight {
		return nil, common.New(constant.InvalidParam, "参数错误")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	maxX := utils.MinInt(global.MapWidth, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)

	rb := make([]model.MapRoleCity, 0)
	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			v, ok := r.posCity[posId]
			if ok {
				rb = append(rb, v.ToModel().(model.MapRoleCity))
			}
		}
	}
	return rb, nil
}

func (r *roleCityService) Get(id int) (*data.MapRoleCity, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	city, ok := r.dbCity[id]
	if ok {
		return city, ok
	}
	return nil, ok
}

func (r *roleCityService) GetMainCity(id int) *data.MapRoleCity {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	cities, ok := r.roleCity[id]
	if ok {
		for _, v := range cities {
			if v.IsMain == 1 {
				return v
			}
		}
	}
	return nil
}

func (r *roleCityService) GetCityCost(id int) int8 {
	return CityFacilityService.GetCost(id) + gameConfig.Basic.City.Cost
}
