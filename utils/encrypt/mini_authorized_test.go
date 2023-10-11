package encrypt

import "testing"

func TestNewMiniAuthorized(t *testing.T) {
	a := NewMiniAuthorized("c8c93222583741bd828579b3d3efd43b_1")
	b := a.Decrypt("Basic MTIzNDU2Nzg5LjE2NTMyOTI4NjUuMTY1Mzg5NzY2NS4yNzQ0YmE3YWQ0MjI3MjYxNDk2YTI0ZTYzYzQ5MDA2ZQ==")
	t.Log(b.Data())
	t.Log(b.Verify())
}
