package main

// Исходный плохой код, который нужно переписать:
func DownloadFile(userID string, fileID string) (string, error) {
    if userID != "" {
        if fileID != "" {
            userExists := checkUser(userID)
            if userExists {
                hasAccess := checkAccess(userID, fileID)
                if hasAccess {
                    return "file_data_bytes", nil
                } else {
                    return "", errors.New("access denied")
                }
            } else {
                return "", errors.New("user not found")
            }
        } else {
            return "", errors.New("empty file id")
        }
    } else {
        return "", errors.New("empty user id")
    }
}

func GoodDownloadFile(userID string, fileID string) (string, error) {
    if userID == "" {
    	return "", errors.New("empty user id")
    }

    if fileID == "" {
        return "", errors.New("empty file id")
    }

    userExists := checkUser(userID)

    if !userExists {
    	return "", errors.New("user not found")
    }

    hasAccess := checkAccess(userID, fileID)

    if !hasAccess { 
    	return "", errors.New("access denied")
    }

    return "file_data_bytes", nil
}