package size

func String(key string, _ [][]byte, length int) int {
	return 56 + len(key) + 8 + length + 1
}

// List 如果是全采样，length = -1；否则 length 为 List 长度，Set、Hash、Zset同理
func List(key string, members [][]byte, length int) int {
	totalMemberSize := 0
	for _, member := range members {
		totalMemberSize += len(member)
	}
	sample := len(members)
	if length > 0 {
		return keySize(key) + length*(1+1+totalMemberSize/sample)
	}
	return keySize(key) + totalMemberSize + (1+1)*sample
}

func Set(key string, members [][]byte, length int) int {
	return List(key, members, length)
}

func Hash(key string, memberValues [][]byte, length int) int {
	var (
		totalMemberSize = 0
		totalValueSize  = 0
	)
	for i := 0; i < len(memberValues); i += 2 {
		totalMemberSize += len(memberValues[i])
		totalValueSize += len(memberValues[i+1])
	}
	sample := len(memberValues) / 2
	if length > 0 {
		return keySize(key) + length*(4+1+totalMemberSize/sample+1+1+totalValueSize/sample)
	}
	return keySize(key) + (4+1)*sample + totalMemberSize + (1+1)*sample + totalValueSize
}

// Zset 按照8字节计算 score 长度
func Zset(key string, members [][]byte, length int) int {
	totalMemberSize := 0
	for _, member := range members {
		totalMemberSize += len(member)
	}
	sample := len(members)
	if length > 0 {
		return keySize(key) + length*(5+totalMemberSize/sample+11+8)
	}
	return keySize(key) + totalMemberSize + 5*sample + (11+8)*sample
}

func keySize(key string) int {
	return 56 + len(key)
}
