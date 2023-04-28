package resp

import "github.com/gin-gonic/gin"

type n struct {
	s
}

func (n n) Next() bool {
	return true
}

func (n n) Send(ctx *gin.Context) {
	payload := n.GetResponsePayload()
	payload.Header.ForEach(func(key, value string) {
		ctx.Header(key, value)
	})
	ctx.Status(payload.Status)
	ctx.Next()
}

func Next() Response {
	return n{
		s: S().(s),
	}
}

// N alias for resp.Next
func N() Response {
	return Next()
}
