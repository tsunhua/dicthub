package dao

import (
	"app/service/dict/model"
	"testing"
)

func TestTraversal(t *testing.T) {
	catalogTree := []*model.TreeNode{
		{
			Id:   "1",
			Name: "標題1",
			Next: []*model.TreeNode{
				{
					Id:   "1-1",
					Name: "標題1-1",
					Next: []*model.TreeNode{
						{
							Id:   "1-1-1",
							Name: "標題1-1-1",
						},
					},
				},
			},
		},
		{
			Id:   "2",
			Name: "標題2",
		},
	}
	catalogCtx := &traversalContext{
		lastLevel: -1,
	}
	arr := traversal(catalogTree, catalogCtx)()
	if arr[1].LinkId != "1"+ID_LINKER+"1-1" || arr[1].LinkName != "標題1"+NAME_LINKER+"標題1-1" {
		t.Fail()
	}
}
