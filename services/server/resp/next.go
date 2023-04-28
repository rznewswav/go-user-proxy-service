package resp

func Next() Response {
	r := S().(response)
	r.next = true
	return r
}

// N alias for resp.Next
func N() Response {
	return Next()
}
