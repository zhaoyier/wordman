package product

import (
	"eventag.cn/wordman/wrapx"
	"fmt"
	"github.com/gin-gonic/gin"
)

func ProductService() {
	fmt.Println("===>>>TT")
	router := gin.Default()
	base := wrapx.New("product")
	base.Register(router, new(Hello))

	router.Run(":8080")
}
