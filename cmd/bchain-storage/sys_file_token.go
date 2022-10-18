package main

import (
	"net/http"

	"github.com/google/uuid"
)

func init() {
	RegisterHandle("/sys/file/token", tokenHandler)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) error {
	auth, ok := authBase(r)
	if !ok {
		return writeMsg(w, 401, "auth failed")
	}
	if auth.User == adminUser {
		return writeMsg(w, 401, "can not use 'admin' to manage files")
	}

	sid := r.FormValue("sid")
	if len(sid) == 0 {
		return writeMsg(w, 403, "params failed")
	}

	if r.Method == "POST" {
		if !_handler.DelayToken(sid) {
			return writeMsg(w, 403, "token has expired")
		}
		return writeMsg(w, 200, "success")
	}

	if r.Method == "DELETE" {
		_handler.DeleteToken(sid)
		return writeMsg(w, 200, "success")
	}

	token := uuid.New().String()
	_handler.AddToken(auth.SpaceName, sid, token)
	return writeMsg(w, 200, token)
}
