package rediskey

const (
	revokeAuthPrefix = "revoke_auth_"
)

func GetRevokeAuthPrefix() string {
	return revokeAuthPrefix
}
