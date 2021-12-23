package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCategory(t *testing.T) {
	text := `# 人物/renwu
## 泛稱/chenghu
### 人稱·指代/rencheng-zhidai
### 一般指稱·尊稱/yibangzhicheng-zuncheng
### 詈稱·貶稱/licheng-biancheng
## 性格/xingge
# 動詞/dongci`
	nodes := parse2TreeNodeBOs(text)

	assert.Equal(t, 7, len(nodes))

	assert.Equal(t, 1, nodes[0].Level)
	assert.Equal(t, "人物", nodes[0].Name)
	assert.Equal(t, "renwu", nodes[0].Id)

	assert.Equal(t, 3, nodes[2].Level)
	assert.Equal(t, "人稱·指代", nodes[2].Name)
	assert.Equal(t, "rencheng-zhidai", nodes[2].Id)
	
	assert.Equal(t, "1.2", nodes[5].Number)
}
