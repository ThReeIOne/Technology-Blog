package lib

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GinCtx(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GqlCtxKeyGin)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	c, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return c, nil
}
