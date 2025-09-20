package captcha

type CheckForm struct {
	Key string `form:"key" binding:"required"`
	X   int    `form:"x" binding:"required"`
	Y   int    `form:"y" binding:"required"`
}
