package controllers

import (
	"fmt"
	"strconv"
	"uu/models"

	"github.com/astaxie/beego"
)

type NovelController struct {
	beego.Controller
}

func (n *NovelController) Get() {
	query := n.Ctx.Request.URL.Query()
	p := query.Get("page")
	if len(p) == 0 {
		p = "0"
	}
	page, err := strconv.Atoi(p)
	fmt.Printf("======  page : %v\n", page)
	if err != nil {
		n.Data["json"] = err
	} else {
		list, err := models.GetNovelsOnPage(page)
		if nil != err {
			n.Data["json"] = err
		} else {
			n.Data["json"] = &list
		}
	}
	n.ServeJSON()
}

func (n *NovelController) getPage() {
	query := n.Ctx.Request.URL.Query()
	p := query.Get("page")
	if len(p) == 0 {
		p = "0"
	}
	page, err := strconv.Atoi(p)
	fmt.Printf("======  page : %v\n", page)
	if err != nil {
		n.Data["json"] = err
	} else {
		list, err := models.GetNovelsOnPage(page)
		if nil != err {
			n.Data["json"] = err
		} else {
			n.Data["json"] = &list
		}
	}
	n.ServeJSON()
}
