package utils

func FindAndDelete(data []string, delete string) []string {
	var respuesta []string
	for i := 0; i < len(data); i++ {
		if data[i] != delete {
			respuesta = append(respuesta, data[i])
		}
	}
	return respuesta
}
