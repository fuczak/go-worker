package handler

import log "github.com/sirupsen/logrus"

func Error(err error) {
	if err != nil {
		log.Error(err.Error())
	}
}
