package gameConfig

import (
	"encoding/json"
	"log"
	"os"
)

type cfg struct {
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	Grain    int    `json:"grain"`
	Wood     int    `json:"wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Durable  int    `json:"durable"`
	Defender int    `json:"defender"`
}

type mapBuildConf struct {
	Title  string `json:"title"`
	Cfg    []cfg  `json:"cfg"`
	cfgMap map[int8][]cfg
}

var MapBuildConf = &mapBuildConf{}

const mapBuildConfFile = "/conf/game/map_build.json"

func (m *mapBuildConf) Load() {
	//获取当前文件路径
	currentDir, _ := os.Getwd()
	//配置文件位置
	cf := currentDir + mapBuildConfFile
	//打包后 程序参数加入配置文件路径
	if len(os.Args) > 1 {
		if path := os.Args[1]; path != "" {
			cf = path + mapBuildConfFile
		}
	}
	data, err := os.ReadFile(cf)
	if err != nil {
		log.Println("地图配置资源读取失败")
		panic(err)
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		log.Println("地图配置资源格式定义失败")
		panic(err)
	}
}
