package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/global"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"Turing-Go/utils"
	"log"
)

var RoleBuildService = &roleBuildService{
	posRB:  make(map[int]*data.MapRoleBuild),
	roleRB: make(map[int][]*data.MapRoleBuild),
}

type roleBuildService struct {
	posRB  map[int]*data.MapRoleBuild
	roleRB map[int][]*data.MapRoleBuild
}

func (r *roleBuildService) Load() {
	total, err := db.Engine.Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).Count(new(data.MapRoleBuild))
	if err != nil {
		log.Println("加载角色建筑异常", err)
		panic(err)
		return
	}
	sysBuild := gameConfig.MapRes.SysBuild
	if int(total) != len(sysBuild) {
		//对不上，需要将系统建筑存入数据库
		//先删除 后插入
		_, err := db.Engine.Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).Delete(new(data.MapRoleBuild))
		if err != nil {
			log.Println("删除角色建筑异常", err)
		}
		for _, v := range sysBuild {
			build := data.MapRoleBuild{
				RId:   0,
				Type:  v.Type,
				Level: v.Level,
				X:     v.X,
				Y:     v.Y,
			}
			build.Init()
			_, err := db.Engine.InsertOne(&build)
			if err != nil {
				log.Println("插入角色建筑异常", err)
			}
		}
	}
	//查找所有的角色建筑
	dbRb := make(map[int]*data.MapRoleBuild)
	err = db.Engine.Find(dbRb)
	if err != nil {
		log.Println("加载角色建筑异常", err)
		panic(err)
	}
	//将其转换为 角色id-建筑 位置-建筑
	for _, v := range dbRb {
		v.Init()
		pos := global.ToPosition(v.X, v.Y)
		r.posRB[pos] = v
		_, ok := r.roleRB[v.RId]
		if !ok {
			r.roleRB[v.RId] = make([]*data.MapRoleBuild, 0)
		} else {
			r.roleRB[v.RId] = append(r.roleRB[v.RId], v)
		}
	}
}

func (r *roleBuildService) GetBuilds(rid int) ([]model.MapRoleBuild, error) {
	mrs := make([]data.MapRoleBuild, 0)
	mr := &data.MapRoleBuild{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	modelBuilds := make([]model.MapRoleBuild, len(mrs))
	if err != nil {
		log.Println("查询角色建筑异常", err)
		return modelBuilds, common.New(constant.DBError, "数据库错误")
	}
	for i, roleBuild := range mrs {
		modelBuilds[i] = roleBuild.ToModel().(model.MapRoleBuild)
	}
	return modelBuilds, nil
}

func (r *roleBuildService) ScanBlock(req *model.ScanBlockReq) ([]model.MapRoleBuild, error) {
	x := req.X
	y := req.Y
	length := req.Length
	var mrbs []model.MapRoleBuild
	if x < 0 || x >= global.MapWidth || y < 0 || y >= global.MapHeight {
		return mrbs, nil
	}

	maxX := utils.MinInt(global.MapWidth, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)

	for i := x - length; i <= maxX; i++ {
		for j := y - length; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			mrb, ok := r.posRB[posId]
			if ok {
				mrbs = append(mrbs, mrb.ToModel().(model.MapRoleBuild))
			}
		}
	}
	return mrbs, nil
}
