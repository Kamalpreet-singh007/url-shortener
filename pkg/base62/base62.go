package base62


import(
	"fmt"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"


func Encode(ID int64)string{
	
	var shortCode string
	for ID > 0 {
		shortCode = string(alphabet[ID % 62]) + shortCode
		ID = ID / 62
	}
	return shortCode
}

func Decode(shortCode string) (int64, error){

	result := int64(0)
	for _, char := range shortCode {
		val:=  int64(strings.Index(alphabet, string(char)))
		if (val == -1) {
			return 0, fmt.Errorf("invalid character in short code: %s", string(char))}
		result = result * 62 +val
	}
	return result, nil
}
