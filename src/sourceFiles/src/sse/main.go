package sse

func Init()  {
	brokerByUser = map[string]broker{}
	//go func() {
	//	cnt := 0
	//	for {
	//		time.Sleep(3 *time.Second)
	//		SendJson("1", map[string]interface{}{"cnt": cnt})
	//		cnt++
	//	}
	//}()
}