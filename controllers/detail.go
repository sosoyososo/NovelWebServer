package controllers

import "github.com/astaxie/beego"
import "uu/models"

type NovelChapterDetailController struct {
	beego.Controller
}

func (n *NovelChapterDetailController) Get() {
	query := n.Ctx.Request.URL.Query()
	url := query.Get("url")
	if len(url) > 0 {
		chapter, err := models.GetChapterDetailWithURL(url)
		if nil == err {
			n.Data["json"] = chapter
		} else {
			n.Data["json"] = err
		}
	} else {
		n.Data["json"] = models.NewModelError(-1, "参数为空")
	}
	n.ServeJSON()
}
