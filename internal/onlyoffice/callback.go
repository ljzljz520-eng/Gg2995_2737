package onlyoffice

import "fmt"

type Callback struct {
	Status                   int
	Key, UserID, DownloadURL string
	Content                  []byte
}
type CallbackResult struct {
	Save    bool
	Close   bool
	Message string
}

func HandleCallback(callback Callback) (CallbackResult, error) {
	if callback.Key == "" || callback.UserID == "" {
		return CallbackResult{}, fmt.Errorf("callback identity is incomplete")
	}
	switch callback.Status {
	case 1:
		return CallbackResult{Message: "editing"}, nil
	case 2:
		if len(callback.Content) == 0 {
			return CallbackResult{}, fmt.Errorf("save callback content is empty")
		}
		return CallbackResult{Save: true, Close: true, Message: "ready to save"}, nil
	case 4:
		return CallbackResult{Close: true, Message: "closed without changes"}, nil
	case 6:
		if len(callback.Content) == 0 {
			return CallbackResult{}, fmt.Errorf("forced save content is empty")
		}
		return CallbackResult{Save: true, Message: "forced save"}, nil
	default:
		return CallbackResult{}, fmt.Errorf("unsupported callback status %d", callback.Status)
	}
}
