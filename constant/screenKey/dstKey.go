package screenKey

func Key(level, clusterName string) string {
	return "DST_" + level + "_" + clusterName
}
