package controllers

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

type ControllerError struct {
	code    int
	message string
}

func (c *ErrorController) Error404() {
	c.Data["json"] = map[string]interface{}{"code": 404, "message": "not found"}
	c.ServeJSON()
}
