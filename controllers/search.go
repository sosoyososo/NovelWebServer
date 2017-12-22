package controllers

import "github.com/astaxie/beego"
import "strings"
import "strconv"
import "uu/models"

type NovelSearchController struct {
	beego.Controller
}

func (n *NovelSearchController) Get() {
	path := n.Ctx.Request.URL.Path
	pathComponents := strings.Split(path, "/")
	typeStr := "title"
	if len(pathComponents) > 2 {
		typeStr = pathComponents[2]
	}

	key := n.Ctx.Request.URL.Query().Get("key")
	pageStr := n.Ctx.Request.URL.Query().Get("page")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		n.Data["json"] = err
	} else {
		var list *[]models.Novel
		var err error
		switch typeStr {
		case "author":
			list, err = models.SearchNovelsByAuthor(key, page)
		case "summary":
			list, err = models.SearchNovelsBySummary(key, page)
		case "tags":
			list, err = models.SearchNovelsByTags(key, page)
		default:
			list, err = models.SearchNovelsByTitle(key, page)
		}
		if err != nil {
			n.Data["json"] = err
		} else {
			n.Data["json"] = list
		}
	}
	n.ServeJSON()
}
