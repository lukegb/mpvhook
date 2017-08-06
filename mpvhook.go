// Copyright 2017 Luke Granger-Brown. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	Scopes = []string{"https://www.googleapis.com/auth/drive.readonly", "email", "profile"}

	OAuthConfig = &oauth2.Config{
		ClientID:     "1032160245384-6pboet8pqv0p409iic5kugvq54v0egtf.apps.googleusercontent.com",
		ClientSecret: "gw7azXzMc26FeNuQGnoFRM0R",
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       Scopes,
		Endpoint:     google.Endpoint,
	}
)

const (
	tokenPath = "/home/lukegb/mpvhook/token.json"
)

func loadToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var t oauth2.Token
	if err := json.NewDecoder(f).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func saveToken(t *oauth2.Token) error {
	fn := tokenPath + ".new"
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(t); err != nil {
		f.Close()
		os.Remove(fn)
		return err
	}

	if err := f.Close(); err != nil {
		os.Remove(fn)
		return err
	}

	if err := os.Rename(fn, tokenPath); err != nil {
		os.Remove(fn)
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()

	tok, err := loadToken()
	if err != nil {
		log.Fatal(err)
	}

	if tok == nil {
		url := OAuthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

		var code string
		fmt.Printf("> ")
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}
		tok, err = OAuthConfig.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}

		if err := saveToken(tok); err != nil {
			log.Fatal(err)
		}
	}

	if !tok.Valid() {
		src := OAuthConfig.TokenSource(ctx, tok)

		var err error
		tok, err = src.Token()
		if err != nil {
			log.Fatal(err)
		}

		if err := saveToken(tok); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(tok.AccessToken)
}
