package logs

import "os"

const maxLogSize = 10 * 1024 * 1024 // 10 MB

func createLogFileIfNotExists(path string) (*os.File, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Если файл не существует, создаем его
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	// check the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() > maxLogSize {
		// if the file size is more than maxLogSize, delete the file
		file.Close()
		if err := os.Remove(path); err != nil {
			return nil, err
		}

		newFile, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		return newFile, nil
	}

	return file, nil
}
