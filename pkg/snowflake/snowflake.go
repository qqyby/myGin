package snowflake

/*
	snowflake 是 twitter 开源的一种分布式 ID 生成算法。该算法可以保证在不借助第三方服务或者说工具的情况下，在多台机器上持续生成保证唯一性的 64 位整型 ID。
*/

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(startTime string, machineID int64) error {
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return err
}

func GenID() string {
	return node.Generate().String()
}