package request

type Register struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=191" zh_field:"姓名" en_field:"name"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required" zh_field:"手機" en_field:"mobile"`
	Password string `form:"password" json:"password" binding:"required" zh_field:"密碼" en_field:"password"`
}

type Login struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required" zh_field:"手機" en_field:"mobile"`
	Password string `form:"password" json:"password" binding:"required" zh_field:"密碼" en_field:"password"`
}
