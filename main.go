package main

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type heroInfo struct {
	Ename int    `json:"ename"`
	Cname string `json:"cname"`
}

var basepath string

func download(url, fileName string) (bool, string) {
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Sprintf(" 下载失败 -- 下载时错误：%s", err.Error())
	}
	defer resp.Body.Close()

	f, err := os.Create(fileName)
	if err != nil {
		return false, fmt.Sprintf(" 下载失败 -- 创建时错误：%s", err.Error())
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return false, fmt.Sprintf(" 下载失败 -- 保存时错误：%s", err.Error())
	}
	return true, " 下载成功"
}

func PathExists(path string) (isExists, isDir bool) {
	pathinfo, err := os.Stat(path)
	if err == nil {
		return true, pathinfo.IsDir()
	}
	return false, false
}

func DirCreate(dir_path string) {
	isExist, isDir := PathExists(dir_path)
	if !isDir {
		if isExist {
			os.Remove(dir_path)
		}
		os.MkdirAll(dir_path, 0755)
	}
}

func saveHeroSkin(info heroInfo) {
	var wg sync.WaitGroup
	heroID := info.Ename
	heroName := info.Cname
	heroPath := path.Join(basepath, heroName)
	fmt.Printf("当前正在下载 %s \t的皮肤 \n", heroName)

	DirCreate(heroPath)
	doc, err := htmlquery.LoadURL(fmt.Sprintf("https://pvp.qq.com/web201605/herodetail/%d.shtml", heroID))
	if err != nil {
		panic(err)
	}
	list := htmlquery.Find(doc, "//div/ul[@class='pic-pf-list pic-pf-list3']/@data-imgname")
	if len(list) == 1 {
		skins := htmlquery.InnerText(list[0])
		skin_names := strings.Split(skins, "|")
		for k, skin_name := range skin_names {
			num := k + 1
			skin_name = strings.Split(skin_name, "&")[0]
			imgurl := fmt.Sprintf("https://game.gtimg.cn/images/yxzj/img201606/skin/hero-info/%d/%d-bigskin-%d.jpg", heroID, heroID, num)
			img_path := path.Join(heroPath, skin_name+".jpg")
			wg.Add(1)
			go func(sk string) {
				defer wg.Done()
				res, msg := download(imgurl, img_path)
				if !res {
					fmt.Println(sk + msg)
				}
			}(skin_name)
		}
		wg.Wait()

	}
}

func main() {
	basepath, _ = os.Getwd()
	basepath = path.Join(basepath, "Heros")
	DirCreate(basepath)

	client := &http.Client{}

	resp, err := client.Get("https://pvp.qq.com/web201605/js/herolist.json")

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var heros []heroInfo
	json.Unmarshal(respBody, &heros)
	for _, hero := range heros {
		saveHeroSkin(hero)
	}

}
