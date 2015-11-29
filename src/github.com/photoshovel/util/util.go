package util

import (
	"os"
)

func CreateOrReopenFileForAppending(path string, fileMode os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, fileMode)
	
	// Handles case of preexisting file.
	if err != nil && os.IsExist(err) {
		var reopenErr error
		file, reopenErr = os.OpenFile(path, os.O_APPEND|os.O_RDWR, fileMode)
		if reopenErr != nil {
			err = reopenErr
		} else {
			err = nil
		}
	} /*else if err != nil && os.IsNotExist(err) {
		// File does not exist.  Weren't able to open, so permissions, capacity, 
		// ??? problems encountered.
		// TODO Handle this better?
	}*/
	
	return
}
