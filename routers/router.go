package routers

import (
	"uu/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/novels", &controllers.NovelController{})
}
