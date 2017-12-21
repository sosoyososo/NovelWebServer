package routers

import (
	"uu/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/novels", &controllers.NovelController{})
	beego.Router("/chapters", &controllers.NovelChaptersController{})
	beego.Router("/detail", &controllers.NovelChapterDetailController{})
}
