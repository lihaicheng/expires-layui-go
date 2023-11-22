package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/dao"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/model"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Menu struct {
	HomeInfo *MenuItem   `json:"homeInfo"`
	LogoInfo *MenuItem   `json:"logoInfo"`
	MenuInfo []*MenuItem `json:"menuInfo"`
}

type MenuItem struct {
	Title  string      `json:"title,omitempty"`
	Image  string      `json:"image,omitempty"`
	Icon   string      `json:"icon,omitempty"`
	Href   string      `json:"href"`
	Target string      `json:"target,omitempty"`
	Child  []*MenuItem `json:"child,omitempty"`
}

func DefaultInitMenuHandler(c *gin.Context) {

	filter := make(map[string]interface{}, 0)
	filter["type"] = "home"
	homeInfo, _ := dao.Menu().Get(filter)
	zap.L().Info("DefaultInitMenuHandler: homeInfo code is " + homeInfo.Code)

	filter["type"] = "logo"
	logoInfo, _ := dao.Menu().Get(filter)
	zap.L().Info("DefaultInitMenuHandler: logoInfo code is " + logoInfo.Code)

	filter["type"] = "head"
	headInfos, _ := dao.Menu().List(filter)
	res := &Menu{
		HomeInfo: &MenuItem{
			Title: homeInfo.Title,
			Href:  homeInfo.Href,
		},
		LogoInfo: &MenuItem{
			Title: logoInfo.Title,
			Image: logoInfo.Image,
			Href:  logoInfo.Href,
		},
		MenuInfo: make([]*MenuItem, 0),
	}

	// 把head先加入
	for i, headInfo := range headInfos {
		zap.L().Info("DefaultInitMenuHandler: menuInfo code is " + headInfo.Code)
		headInfoChild := strings.Split(headInfo.Child, model.SplitFlag4Child)
		if headInfo.IsShown == 0 {
			continue
		}
		res.MenuInfo = append(res.MenuInfo, &MenuItem{
			Title:  headInfo.Title,
			Icon:   headInfo.Icon,
			Href:   headInfo.Href,
			Target: headInfo.Target,
			Child:  make([]*MenuItem, 0),
		})
		// 把directory和menu加入
		for j, menuCode := range headInfoChild {
			filter = make(map[string]interface{}, 0)
			filter["code"] = menuCode
			menuInfo, _ := dao.Menu().Get(filter)
			//menuInfoChild := strings.Split(headInfo.Child, model.SplitFlag4Child)
			if menuInfo.IsShown == 0 {
				continue
			}
			if menuInfo.Type == "menu" {
				res.MenuInfo[i].Child = append(res.MenuInfo[i].Child, &MenuItem{
					Title:  menuInfo.Title,
					Href:   menuInfo.Href,
					Icon:   menuInfo.Icon,
					Target: menuInfo.Target,
				})
			} else {
				menuInfoChild := strings.Split(menuInfo.Child, model.SplitFlag4Child)
				res.MenuInfo[i].Child = append(res.MenuInfo[i].Child, &MenuItem{
					Title:  menuInfo.Title,
					Href:   menuInfo.Href,
					Icon:   menuInfo.Icon,
					Target: menuInfo.Target,
					Child:  make([]*MenuItem, 0),
				})
				for _, childCode := range menuInfoChild {
					filter = make(map[string]interface{}, 0)
					filter["code"] = childCode
					childInfo, _ := dao.Menu().Get(filter)
					if childInfo.Type == "menu" && childInfo.IsShown == 1 {
						res.MenuInfo[i].Child[j].Child = append(res.MenuInfo[i].Child[j].Child, &MenuItem{
							Title:  childInfo.Title,
							Href:   childInfo.Href,
							Icon:   childInfo.Icon,
							Target: childInfo.Target,
						})
					}
				}
			}
		}
	}

	// 使用c.Data将JSON文本发送到前端
	c.JSON(http.StatusOK, res)

}
