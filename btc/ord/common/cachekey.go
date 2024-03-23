package common

const (
	TapElement = "tap:element:%s"
)

func CacheKey(key string) string {
	//value, exists := os.LookupEnv(EnvTag)
	//if exists {
	//	return fmt.Sprintf("%s:%s", key, value)
	//} else {
	//	return key
	//}
	return key
}
