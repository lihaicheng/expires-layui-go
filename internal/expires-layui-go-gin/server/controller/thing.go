package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/dao"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/model"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Thing struct {
	HomeInfo  *ThingItem   `json:"homeInfo"`
	LogoInfo  *ThingItem   `json:"logoInfo"`
	ThingInfo []*ThingItem `json:"thingInfo"`
}

type ThingItem struct {
	Title  string       `json:"title,omitempty"`
	Image  string       `json:"image,omitempty"`
	Icon   string       `json:"icon,omitempty"`
	Href   string       `json:"href"`
	Target string       `json:"target,omitempty"`
	Child  []*ThingItem `json:"child,omitempty"`
}

func DefaultInitThingHandler(c *gin.Context) {

	filter := make(map[string]interface{}, 0)
	filter["type"] = "home"
	homeInfo, _ := dao.Thing().Get(filter)
	zap.L().Info("DefaultInitThingHandler: homeInfo code is " + homeInfo.Code)

	filter["type"] = "logo"
	logoInfo, _ := dao.Thing().Get(filter)
	zap.L().Info("DefaultInitThingHandler: logoInfo code is " + logoInfo.Code)

	filter["type"] = "head"
	headInfos, _ := dao.Thing().List(filter)
	res := &Thing{
		HomeInfo: &ThingItem{
			Title: homeInfo.Title,
			Href:  homeInfo.Href,
		},
		LogoInfo: &ThingItem{
			Title: logoInfo.Title,
			Image: logoInfo.Image,
			Href:  logoInfo.Href,
		},
		ThingInfo: make([]*ThingItem, 0),
	}

	// 把head先加入
	for i, headInfo := range headInfos {
		zap.L().Info("DefaultInitThingHandler: thingInfo code is " + headInfo.Code)
		headInfoChild := strings.Split(headInfo.Child, model.SplitFlag4Child)
		if headInfo.IsShown == 0 {
			continue
		}
		res.ThingInfo = append(res.ThingInfo, &ThingItem{
			Title:  headInfo.Title,
			Icon:   headInfo.Icon,
			Href:   headInfo.Href,
			Target: headInfo.Target,
			Child:  make([]*ThingItem, 0),
		})
		// 把directory和thing加入
		for j, thingCode := range headInfoChild {
			filter = make(map[string]interface{}, 0)
			filter["code"] = thingCode
			thingInfo, _ := dao.Thing().Get(filter)
			//thingInfoChild := strings.Split(headInfo.Child, model.SplitFlag4Child)
			if thingInfo.IsShown == 0 {
				continue
			}
			if thingInfo.Type == "thing" {
				res.ThingInfo[i].Child = append(res.ThingInfo[i].Child, &ThingItem{
					Title:  thingInfo.Title,
					Href:   thingInfo.Href,
					Icon:   thingInfo.Icon,
					Target: thingInfo.Target,
				})
			} else {
				thingInfoChild := strings.Split(thingInfo.Child, model.SplitFlag4Child)
				res.ThingInfo[i].Child = append(res.ThingInfo[i].Child, &ThingItem{
					Title:  thingInfo.Title,
					Href:   thingInfo.Href,
					Icon:   thingInfo.Icon,
					Target: thingInfo.Target,
					Child:  make([]*ThingItem, 0),
				})
				for _, childCode := range thingInfoChild {
					filter = make(map[string]interface{}, 0)
					filter["code"] = childCode
					childInfo, _ := dao.Thing().Get(filter)
					if childInfo.Type == "thing" && childInfo.IsShown == 1 {
						res.ThingInfo[i].Child[j].Child = append(res.ThingInfo[i].Child[j].Child, &ThingItem{
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
