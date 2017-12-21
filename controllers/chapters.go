package controllers

import "github.com/astaxie/beego"
import "uu/models"
import "strconv"

type NovelChaptersController struct {
	beego.Controller
}

func (n *NovelChaptersController) Get() {
	query := n.Ctx.Request.URL.Query()
	id := query.Get("id")
	page := query.Get("page")
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 0
	}

	if len(id) == 0 {
		n.Data["json"] = models.NullParameter
	} else {
		list, err := models.GetChaptersForNovelWithId(id, p)
		if err != nil {
			n.Data["json"] = err
		} else if nil != list {
			n.Data["json"] = &list
		} else {
			n.Data["json"] = &[]string{}
		}
	}
	n.ServeJSON()
}
