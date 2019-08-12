package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

// ErrNoAvatar is returned when there is no avatar could be found
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents type capable to serve user profile pictures
type Avatar interface {
	GetAvatarURL(user ChatUser) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(user ChatUser) (string, error) {
	url := user.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(user ChatUser) (string, error) {
	return fmt.Sprintf("//www.gravatar.com/avatar/%s", user.UniqueID()), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(user ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := path.Match(user.UniqueID()+"*", file.Name()); match {
			return "/avatars/" + file.Name(), nil
		}
	}
	return "", ErrNoAvatarURL
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(user ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(user); err == nil {
			return url, nil
		}
	}

	return "", ErrNoAvatarURL
}
