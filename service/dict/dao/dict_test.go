package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCategory(t *testing.T) {
	text := "# 人物/renwu\r\n## 泛稱/chenghu/http://www.baidu.com\r\n### 人稱·指代/rencheng-zhidai\r\n### 一般指稱·尊稱/yibangzhicheng-zuncheng\r\n### 詈稱·貶稱/licheng-biancheng\r\n## 性格/xingge\r\n# 動詞/dongci\r\n## 肢體動作/zhitidongzuo"
	nodes := parse2TreeNodeBOs(text)

	assert.Equal(t, 8, len(nodes))

	assert.Equal(t, 1, nodes[0].Level)
	assert.Equal(t, "人物", nodes[0].Name)
	assert.Equal(t, "renwu", nodes[0].Id)

	assert.Equal(t, 3, nodes[2].Level)
	assert.Equal(t, "人稱·指代", nodes[2].Name)
	assert.Equal(t, "rencheng-zhidai", nodes[2].Id)
	assert.Equal(t, "renwu.chenghu.rencheng-zhidai", nodes[2].LinkId)
	assert.Equal(t, "人物/泛稱/人稱·指代", nodes[2].LinkName)
	
	assert.Equal(t, "1.2", nodes[5].Number)
	assert.Equal(t, "2.1", nodes[7].Number)
}
