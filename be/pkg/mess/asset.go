package mess

type Assets map[AssetKey]AssetData
type AssetKey string
type AssetData []byte

func NewAssetKey(str string) AssetKey {
	if str[:1] != "/" {
		return AssetKey("/" + str)
	}
	return AssetKey(str)
}
