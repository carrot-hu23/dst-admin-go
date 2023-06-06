package dstKey

func GetMasterKey(clusterName string) string {
	return "DST_" + clusterName + "_Master"
}

func GetCavesKey(clusterName string) string {
	return "DST_" + clusterName + "_Caves"
}
