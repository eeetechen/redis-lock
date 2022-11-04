package distributed_lock

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func GenerateRedisLockVal(rl *RedisLock) string {
	strStatus := strconv.FormatInt(int64(rl.status), 10)
	strTime := strconv.FormatInt(rl.timestamp, 10)
	strRand := strconv.FormatInt(int64(rl.randVal), 10)
	key := fmt.Sprintf(strStatus + "." + strTime + "." + strRand)
	return key
}

func ParseRedisLockVal(result string) (uint8, int64, int, error) {
	valSli := strings.Split(result, ".")
	if len(valSli) != 3 {
		zap.S().Error("parse distributed_lock failure")
		return 0, 0, 0, fmt.Errorf("parse distributed_lock failure len(result) != 2 ")
	}
	Status, err := strconv.Atoi(valSli[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parse status err val :%v", valSli[0])
	}
	Timestamp, err := strconv.Atoi(valSli[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parse Timestamp err val :%v", valSli[1])
	}
	RandVal, err := strconv.Atoi(valSli[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parse randVal err val :%v", valSli[1])
	}
	return uint8(Status), int64(Timestamp), int(RandVal), nil
}
